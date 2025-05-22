.PHONY: run build test clean docker-build docker-run docker-stop swagger lint help all docker-compose-up docker-compose-down

# Binary name
BINARY_NAME=nwc_app

# Go source files
GO_FILES=$(shell find . -name "*.go" -type f)

# Current version
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

# Main package path
MAIN_PACKAGE=./

# Default target
all: lint test build

# Run the application
run:
	@echo "Running application..."
	@go run .

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) $(MAIN_PACKAGE)

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(BINARY_NAME)
	@go clean

# Build docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(BINARY_NAME):$(VERSION) .

# Run docker container
docker-run:
	@echo "Running Docker container..."
	@docker run -p 8080:8080 --env-file .env --name $(BINARY_NAME) $(BINARY_NAME):$(VERSION)

# Stop docker container
docker-stop:
	@echo "Stopping Docker container..."
	@docker stop $(BINARY_NAME)
	@docker rm $(BINARY_NAME)

# Docker compose up
docker-up:
	@echo "Starting with Docker Compose..."
	@docker-compose up -d

# Docker compose down
docker-down:
	@echo "Stopping with Docker Compose..."
	@docker-compose down

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	@go run github.com/swaggo/swag/cmd/swag init

# Lint the code
lint:
	@echo "Linting code..."
	@if command -v golangci-lint >/dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Please install it by running:"; \
		echo "go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Display help information
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  run              Run the application locally"
	@echo "  build            Build the application"
	@echo "  test             Run tests"
	@echo "  clean            Clean build artifacts"
	@echo "  docker-build     Build Docker image"
	@echo "  docker-run       Run Docker container"
	@echo "  docker-stop      Stop Docker container"
	@echo "  docker-compose-up Start with Docker Compose"
	@echo "  docker-compose-down Stop with Docker Compose"
	@echo "  swagger          Generate Swagger documentation"
	@echo "  lint             Lint the code"
	@echo "  help             Display this help message"
	@echo "  all              Run lint, test, and build"
