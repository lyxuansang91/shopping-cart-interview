FROM golang:1.24.1-alpine as builder

# Install system dependencies
RUN apk add --no-cache git

# Install Air
RUN go install github.com/air-verse/air@v1.61.7

# Set working directory
WORKDIR /app

# Copy the entire project structure to maintain relative paths
COPY . .

# Set working directory for the service
WORKDIR /app/services/main

# Download dependencies
RUN go mod download

# Build the application
RUN go build -o /app/bin/main ./cmd/main.go

# Final stage
FROM golang:1.24.1-alpine

# Install system dependencies
RUN apk add --no-cache git

# Install Air
RUN go install github.com/air-verse/air@v1.61.7

# Set working directory
WORKDIR /app/services/main

# Copy the binary from builder stage
COPY --from=builder /app/bin/main /app/bin/main

# Copy the source code for Air
COPY services/main/ .

# Use CMD to handle the USE_AIR environment variable
CMD ["sh", "-c", "if [ \"$USE_AIR\" = \"backend-service\" ]; then echo \"Starting with Air for hot reloading...\" && air -c .air.toml; else echo \"Starting without Air...\" && /app/bin/main; fi"]
