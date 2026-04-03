app_name := "app"
build_dir := "build"

# Build the main binary
build:
    @echo "Building {{app_name}}..."
    @mkdir -p {{build_dir}}
    go build -o {{build_dir}}/{{app_name}} ./cmd

# Run a specific service: just run api-gateway, just run ws-server, etc.
run service:
    go run ./cmd {{service}}

# Build for Linux amd64
build-linux:
    @mkdir -p {{build_dir}}
    GOOS=linux GOARCH=amd64 go build -o {{build_dir}}/{{app_name}}-linux-amd64 ./cmd

# Build for macOS
build-darwin:
    @mkdir -p {{build_dir}}
    GOOS=darwin GOARCH=amd64 go build -o {{build_dir}}/{{app_name}}-darwin-amd64 ./cmd
    GOOS=darwin GOARCH=arm64 go build -o {{build_dir}}/{{app_name}}-darwin-arm64 ./cmd

# Clean build artifacts
clean:
    rm -rf {{build_dir}}

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

# Generate code
generate:
    go generate ./...

# Start docker services
docker-up:
    docker-compose up -d

# Stop docker services
docker-down:
    docker-compose down

# Tail docker logs
docker-logs:
    docker-compose logs -f

# Install dev tools
install:
    go install golang.org/x/tools/cmd/goimports@latest
