# Variables
PROJECT_NAME = url-shortener
DOCKER_IMAGE = $(PROJECT_NAME):latest
DOCKER_CONTAINER = $(PROJECT_NAME)-container
PORT = 8080

# Default target
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make setup    - Setup the project (install dependencies)"
	@echo "  make up       - Build and run Docker container"
	@echo "  make down     - Stop and remove Docker container"
	@echo "  make test     - Run tests"
	@echo "  make clean    - Clean up build artifacts"

# Setup the project
.PHONY: setup
setup:
	@echo "Setting up the project..."
	go mod tidy
	@echo "Project setup complete!"

# Build and run Docker container
.PHONY: up
up:
	@echo "Building and starting Docker container..."
	docker build -t $(DOCKER_IMAGE) .
	docker run -d --name $(DOCKER_CONTAINER) -p $(PORT):8080 $(DOCKER_IMAGE)
	@echo "Container started! Service available at http://localhost:$(PORT)"

# Stop and remove Docker container
.PHONY: down
down:
	@echo "Stopping and removing Docker container..."
	-docker stop $(DOCKER_CONTAINER)
	-docker rm $(DOCKER_CONTAINER)
	@echo "Container stopped and removed!"

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test ./... -v
	@echo "Tests completed!"

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test ./... -v -cover
	@echo "Tests with coverage completed!"

# Run tests with race detection
.PHONY: test-race
test-race:
	@echo "Running tests with race detection..."
	go test ./... -v -race
	@echo "Tests with race detection completed!"

# Clean up build artifacts
.PHONY: clean
clean:
	@echo "Cleaning up build artifacts..."
	-docker rmi $(DOCKER_IMAGE)
	@echo "Cleanup completed!"

# Restart the service (down + up)
.PHONY: restart
restart: down up

# Show logs
.PHONY: logs
logs:
	@echo "Showing container logs..."
	docker logs -f $(DOCKER_CONTAINER)

# Show container status
.PHONY: status
status:
	@echo "Container status:"
	docker ps -a --filter name=$(DOCKER_CONTAINER)

# Run the service locally (without Docker)
.PHONY: run
run:
	@echo "Running service locally..."
	go run main.go

# Install dependencies only
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go mod download
	@echo "Dependencies installed!" 