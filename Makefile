.PHONY: default up down dev logs db clean fmt proto sqlc temporal-setup monitoring dbtool dbtool-migrate dbtool-seed test test-parallel pre-commit-install pre-commit-run fmt-proto

GATEWAY_MOD  := $(shell go list -m -f "{{.Dir}}" github.com/grpc-ecosystem/grpc-gateway/v2)

# Default development shell
default:
	devbox shell

# Start all services in background
up: down temporal-setup monitoring-setup
	@echo "Starting infrastructure services..."
	@docker compose up -d mysql-db redis-cache temporal temporal-admin-tools temporal-ui prometheus grafana loki promtail
# @docker compose up -d mysql-db redis-cache temporal temporal-admin-tools temporal-ui prometheus grafana loki promtail frontend
	@echo "Waiting for services to be ready..."
	@sleep 10  # Give services time to start
	@echo "Setting up Temporal namespace..."
	@docker compose exec -T temporal-admin-tools tctl --ns default namespace describe > /dev/null 2>&1 || \
		docker compose exec -T temporal-admin-tools tctl --ns default namespace register > /dev/null 2>&1
	@echo "Running database migrations and seeders..."
	@make dbtool
	@make dbtool-migrate
	@sleep 5
	@make dbtool-seed
	@echo "Starting Go services with Air..."
	@GOWORK=off ./scripts/start-services.sh

# Setup monitoring requirements
monitoring-setup:
	@echo "Setting up monitoring requirements..."
	@mkdir -p config/prometheus
	@mkdir -p config/grafana/provisioning/datasources
	@mkdir -p config/grafana/provisioning/dashboards
	@mkdir -p config/promtail

# Setup Temporal requirements
temporal-setup:
	@echo "Setting up Temporal requirements..."
	@mkdir -p config/dynamicconfig
	@chmod +x ./scripts/init-temporal-db.sh

# Development mode with live reload and logs
dev: down
	@echo "Starting in development mode with live reload..."
	@GOWORK=off ./scripts/start-services.sh

# Stop and remove all services
down:
	@echo "Stopping all services..."
	@docker-compose down -v
	@echo "Killing any running Air and main processes..."
	@pkill -f "air -c" || true
	@echo "Waiting for ports to be released..."
	@sleep 3
	@echo "Killing processes on configured HTTP and gRPC ports..."
	@for service in services/*; do \
		if [ -f "$$service/.env.local" ]; then \
			http_port=$$(grep "HTTP_PORT=" $$service/.env.local | cut -d'=' -f2); \
			grpc_port=$$(grep "GRPC_PORT=" $$service/.env.local | cut -d'=' -f2); \
			if [ ! -z "$$http_port" ]; then \
				lsof -ti:$$http_port | xargs kill -9 2>/dev/null || true; \
			fi; \
			if [ ! -z "$$grpc_port" ]; then \
				lsof -ti:$$grpc_port | xargs kill -9 2>/dev/null || true; \
			fi; \
		fi; \
	done
	@echo "All services stopped and ports released successfully."

# View logs (usage: make logs service=payments-service)
logs:
	docker compose logs -f $(service)

# Start only the database
db:
	docker compose up -d mysql-db

# Clean temporary files
clean:
	rm -rf services/*/tmp/*
	docker-compose down -v
	docker volume rm go-mod-cache || true

# Format all Go code
fmt:
	@echo "Formatting Go code..."
	@find . -type f -name "go.mod" -exec dirname {} \; | while read dir; do \
		echo "Running go mod tidy and formatting in $$dir"; \
		(cd $$dir && GOWORK=off go mod tidy && go fmt ./...); \
	done
	@which goimports >/dev/null 2>&1 || go install golang.org/x/tools/cmd/goimports@latest
	@goimports -w ./services ./packages

# Generate gRPC and protobuf Go files for all services

proto:
	@echo "Cleaning up existing generated files..."
	@rm -rf packages/proto/pkg/proto/*
	@make fmt-proto
	@echo "Generating gRPC and protobuf Go files..."
	@cd packages/proto && \
	buf dep update && \
	buf generate
	@echo "Proto files generated successfully"
	@echo "Creating softlink for api.swagger.json (for air hot reload, do not commit this file)..."
	@ln -sf ../../packages/proto/pkg/proto/docs/swagger/api.swagger.json services/adapterstripe/api.swagger.json
	@make fmt


# Generate sqlc code for all services
sqlc:
	@echo "Cleaning up existing generated files..."
	@find services -name "sqlc.yaml" -exec dirname {} \; | while read dir; do \
		echo "Cleaning generated files in $$dir/internal/pkg/db"; \
		rm -rf $$dir/internal/pkg/db/*; \
	done
	@echo "Generating sqlc code for all services..."
	@find services -name "sqlc.yaml" -exec dirname {} \; | while read dir; do \
		echo "Running sqlc generate in $$dir"; \
		(cd $$dir && GOWORK=off sqlc generate 2>&1 | grep -v "no queries contained in paths" || true); \
	done

# Parallel run tests for all packages and services
test-parallel:
	@echo "Running tests for all packages and services in parallel..."
	@find . -mindepth 1 -type f -name "go.mod" -exec dirname {} \; | \
		grep -v '^.$$' | \
		xargs -I{} -P 8 sh -c '\
			echo "→ Running tests in {}"; \
			cd "{}"; \
			if go list ./... | grep -q .; then \
				GOWORK=off go test -v -cover ./...; \
			else \
				echo "⚡ No tests to run in {}"; \
			fi \
		'

# Run tests for all packages and services
test:
	@echo "Running tests for all packages and services..."
	@set -e; \
	for dir in $$(find . -mindepth 1 -type f -name "go.mod" -exec dirname {} \;); do \
		echo "→ Running tests in $$dir"; \
		cd $$dir; \
		if go list ./... | grep -q .; then \
			GOWORK=off go test -v -cover ./...; \
		else \
			echo "⚡ No tests to run in $$dir"; \
		fi; \
		cd - >/dev/null; \
	done

# Install pre-commit hooks
pre-commit-install:
	@echo "Installing pre-commit hooks..."
	@pre-commit install

# Run pre-commit hooks on all files
pre-commit-run:
	make pre-commit-install
	@echo "Running pre-commit hooks..."
	@pre-commit run --all-files

# The following are internal targets; Do not use them directly --------------

# Helper target to wait for database
wait-for-db:
	@echo "Waiting for database to be ready..."
	@until docker compose exec mysql-db mysqladmin ping -h localhost -u root -proot --silent; do \
		echo "Waiting for MySQL..."; \
		sleep 2; \
	done

# Add these new targets for Temporal convenience
temporal-ui:
	@echo "Opening Temporal UI in browser..."
	@open http://localhost:8088

temporal-cli:
	@echo "Running Temporal CLI..."
	docker compose exec temporal-admin-tools tctl

temporal-namespace:
	@echo "Creating default namespace..."
	docker compose exec temporal-admin-tools tctl --ns default namespace register

# Monitoring convenience targets
grafana:
	@echo "Opening Grafana in browser..."
	@open http://localhost:3000

prometheus:
	@echo "Opening Prometheus in browser..."
	@open http://localhost:9090

monitoring-logs:
	@echo "Viewing monitoring stack logs..."
	docker compose logs -f prometheus grafana loki

# Database tool targets
dbtool:
	@echo "Building database tool..."
	@cd packages/dbtool && GOWORK=off go build -o ../../bin/dbtool ./cmd/dbtool

dbtool-migrate:
	@echo "Running migrations for all services..."
	@bin/dbtool migrate up --uri "mysql://cinch:cinch@tcp(localhost:3306)/adapterstripe" --services adapterstripe

dbtool-seed:
	@echo "Running seeders for all services..."
	@bin/dbtool seed --uri "mysql://cinch:cinch@tcp(localhost:3306)/adapterstripe" --services adapterstripe

# Generate Swagger documentation from proto files
swagger:
	@echo "Generating Swagger documentation..."
	@cd packages/proto && \
		buf dep update && \
		buf generate

fmt-proto:
	@echo "Formatting proto files..."
	@cd packages/proto && \
		buf format -w .
