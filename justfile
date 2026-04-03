# Build all services
build-all:
    go build -o demo-api/bin/demo-api                   ./demo-api/cmd
    go build -o demo-grpc/bin/demo-grpc                 ./demo-grpc/cmd
    go build -o demo-ws/bin/demo-ws                     ./demo-ws/cmd
    go build -o demo-worker-kafka/bin/demo-worker-kafka ./demo-worker-kafka/cmd
    go build -o demo-worker-order/bin/demo-worker-order ./demo-worker-order/cmd

# Build a specific service: just build demo-api
build service:
    go build -o {{service}}/bin/{{service}} ./{{service}}/cmd

# Run a specific service: just run demo-api
run service:
    go run ./{{service}}/cmd

# Run a specific service with custom config: just run-cfg demo-api ./demo-api/config.development.toml
run-cfg service config:
    go run ./{{service}}/cmd -config {{config}}

# --- Scripts ---

# Run a script command: just script migrate, just script seed, just script fix-order
script +args:
    go run ./scripts/cmd {{args}}

# Build scripts binary
build-scripts:
    go build -o scripts/bin/scripts ./scripts/cmd

# --- Build ---

# Build all services for Linux amd64
build-linux:
    GOOS=linux GOARCH=amd64 go build -o demo-api/bin/demo-api-linux-amd64                   ./demo-api/cmd
    GOOS=linux GOARCH=amd64 go build -o demo-grpc/bin/demo-grpc-linux-amd64                 ./demo-grpc/cmd
    GOOS=linux GOARCH=amd64 go build -o demo-ws/bin/demo-ws-linux-amd64                     ./demo-ws/cmd
    GOOS=linux GOARCH=amd64 go build -o demo-worker-kafka/bin/demo-worker-kafka-linux-amd64 ./demo-worker-kafka/cmd
    GOOS=linux GOARCH=amd64 go build -o demo-worker-order/bin/demo-worker-order-linux-amd64 ./demo-worker-order/cmd

# Clean build artifacts
clean:
    rm -rf */bin

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
