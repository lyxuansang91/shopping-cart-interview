# URL Shortener Service

A simple URL shortener service built with [Echo](https://echo.labstack.com/) in Go.

## Features

- Generate short URLs from long URLs
- Redirect short URLs to original long URLs
- Thread-safe in-memory storage
- RESTful API
- Comprehensive unit tests
- Docker containerization
- Makefile automation
- Environment-based configuration

## Configuration

The service can be configured using environment variables. Copy `env.example` to `.env` and modify as needed:

```bash
cp env.example .env
```

### Environment Variables

| Variable      | Default                 | Description                              |
| ------------- | ----------------------- | ---------------------------------------- |
| `PORT`        | `8080`                  | Port number for the server               |
| `BASE_URL`    | `http://localhost:8080` | Base URL for generating short links      |
| `LOG_LEVEL`   | `info`                  | Logging level (debug, info, warn, error) |
| `ENABLE_CORS` | `true`                  | Enable CORS middleware (true/false)      |

### Example Configuration

```bash
# Production settings
PORT=3000
BASE_URL=https://shortener.example.com
LOG_LEVEL=info
ENABLE_CORS=true

# Development settings
PORT=8080
BASE_URL=http://localhost:8080
LOG_LEVEL=debug
ENABLE_CORS=true
```

## API Endpoints

### 1. Shorten URL

**POST** `/api/shortlinks`

Request Body:

```json
{
  "long_url": "https://www.example.com/very/long/url"
}
```

Response:

```json
{
  "short_url": "http://localhost:8080/shortlinks/abc123",
  "id": "abc123"
}
```

### 2. Redirect to Long URL

**GET** `/api/shortlinks/{id}`
Retrieve details of a short link

```json
{
  "id": "abc123",
  "original_url": "https://example.com",
  "created_at": "2024-01-01T12:00:00Z"
}
```

**GET** `/shortlinks/{id}`

Public redirect endpoint – 302 redirect to the original URL Response: HTTP 302 with Location: original_url

## Quick Start with Makefile

The project includes a comprehensive Makefile for easy development and deployment.

### Prerequisites

- Go 1.24.1 or later
- Docker (for containerized deployment)
- Make

### Initial Setup

```bash
# Setup the project (install dependencies)
make setup
```

### Running the Service

#### Option 1: Docker (Recommended)

```bash
# Build and start the service in Docker
make up

# The service will be available at http://localhost:8080
```

#### Option 2: Local Development

```bash
# Run the service locally (without Docker)
make run
```

### Stopping the Service

```bash
# Stop and remove Docker container
make down
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run tests with race detection
make test-race
```

## Complete Makefile Commands

### Core Commands

```bash
make setup    # Setup the project (install dependencies)
make up       # Build and run Docker container
make down     # Stop and remove Docker container
make test     # Run tests
make clean    # Clean up build artifacts
```

### Additional Commands

```bash
make help           # Show all available commands
make test-coverage  # Run tests with coverage report
make test-race      # Run tests with race detection
make restart        # Restart the service (down + up)
make logs           # Show container logs
make status         # Show container status
make run            # Run locally without Docker
make deps           # Install dependencies only
```

## Development Workflow

### 1. First Time Setup

```bash
make setup
```

### 2. Development Cycle

```bash
# Start the service
make up

# Run tests
make test

# View logs
make logs

# Make changes to your code...

# Restart with changes
make restart

# Stop when done
make down
```

### 3. Testing

```bash
# Run basic tests
make test

# Run tests with coverage
make test-coverage

# Run tests with race detection (for concurrency issues)
make test-race
```

## Manual Setup (Alternative)

If you prefer not to use the Makefile:

### Install dependencies

```bash
go mod tidy
cd services && go mod tidy
```

### Run the service

```bash
go run main.go
```

### Run tests

```bash
go test ./services -v
```

The service will start on the port specified in the `PORT` environment variable (default: 8080)

## Example Usage

**Shorten a URL:**

```bash
curl -X POST http://localhost:8080/api/shortlinks \
  -H "Content-Type: application/json" \
  -d '{"long_url": "https://www.google.com"}'
```

**Access the shortened URL:**

```bash
curl -I http://localhost:8080/shortlinks/abc123
```

**Get short link details:**

```bash
curl http://localhost:8080/api/shortlinks/abc123
```

**Health check:**

```bash
curl http://localhost:8080/health
```

## Docker Commands

If you prefer to use Docker directly:

```bash
# Build the image
docker build -t url-shortener .

# Run the container with custom environment variables
docker run -d --name url-shortener-container \
  -p 8080:8080 \
  -e PORT=8080 \
  -e BASE_URL=http://localhost:8080 \
  url-shortener

# Stop the container
docker stop url-shortener-container

# Remove the container
docker rm url-shortener-container
```

## Notes

- This service uses in-memory storage. All data will be lost when the service restarts.
- For production, consider using a persistent database and adding authentication, rate limiting, and HTTPS.
- The Docker image uses Alpine Linux for a smaller footprint and better security.
- The service runs as a non-root user inside the container for security.
- Environment variables can be set in a `.env` file or passed directly to the container.
