# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /build

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Copy binary
COPY --from=builder /build/app .

# Copy config files
COPY --from=builder /build/demo-api/config.toml ./demo-api/
COPY --from=builder /build/demo-grpc/config.toml ./demo-grpc/
COPY --from=builder /build/demo-ws/config.toml ./demo-ws/
COPY --from=builder /build/demo-worker-kafka/config.toml ./demo-worker-kafka/
COPY --from=builder /build/demo-worker-order/config.toml ./demo-worker-order/

ENTRYPOINT ["./app"]
