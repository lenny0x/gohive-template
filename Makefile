.PHONY: build run-% clean test lint deps

APP_NAME := app
BUILD_DIR := build

# Build the main binary
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd

# Run a specific service: make run-api-gateway, make run-ws-server, etc.
run-%:
	go run ./cmd $*

# Build for specific platforms
build-linux:
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 ./cmd

build-darwin:
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 ./cmd
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 ./cmd

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Lint
lint:
	golangci-lint run ./...

# Install dependencies
deps:
	go mod download
	go mod tidy

# Generate code (if needed)
generate:
	go generate ./...

# Docker compose commands
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

# Help
help:
	@echo "Available commands:"
	@echo "  make build          - Build the application"
	@echo "  make run-<service>  - Run a specific service (e.g., make run-api-gateway)"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make test           - Run tests"
	@echo "  make lint           - Run linter"
	@echo "  make deps           - Install dependencies"
	@echo "  make docker-up      - Start docker services"
	@echo "  make docker-down    - Stop docker services"

install:
	go install golang.org/x/tools/cmd/goimports@latest