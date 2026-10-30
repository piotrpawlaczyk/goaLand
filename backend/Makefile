# Simple Makefile for a Go

# Build the application
build:
	@echo "Building..."
	@mkdir -p bin
	@go build -o bin/main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go $(ARGS)

# Create missing tables and fetch initial data
init-data:
	@go run cmd/api/main.go --init-data

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

# Run unit and integration tests
all-test:
	@echo "Running all tests..."
	@go test -p 1 -tags=integration ./... -v

# Integrations Tests for the application
itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -rf bin

# Live Reload
watch:
	@air

.PHONY: build run init-data test all-test clean watch docker-up docker-down itest
