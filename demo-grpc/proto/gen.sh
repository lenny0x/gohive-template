#!/bin/bash
# Generate Go code from proto files
# Requires: protoc, protoc-gen-go, protoc-gen-go-grpc
#
# Install:
#   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
#   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

SCRIPT_DIR=$(dirname "$0")

protoc --proto_path="$SCRIPT_DIR" \
  --go_out="$SCRIPT_DIR" --go_opt=paths=source_relative \
  --go-grpc_out="$SCRIPT_DIR" --go-grpc_opt=paths=source_relative \
  "$SCRIPT_DIR"/*.proto
