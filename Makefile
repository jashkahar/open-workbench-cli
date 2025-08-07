# Makefile for Open Workbench Platform
# Provides common development tasks

.PHONY: help build test clean install deps lint format

# Default target
help:
	@echo "Open Workbench Platform - Available commands:"
	@echo ""
	@echo "Development:"
	@echo "  build     - Build the binary for current platform"
	@echo "  build-all - Build binaries for all platforms"
	@echo "  test      - Run tests"
	@echo "  clean     - Clean build artifacts"
	@echo "  deps      - Install/update dependencies"
	@echo ""
	@echo "Code Quality:"
	@echo "  lint      - Run linters"
	@echo "  format    - Format code"
	@echo ""
	@echo "Installation:"
	@echo "  install   - Install binary to system PATH"

# Build for current platform
build:
	@echo "Building for current platform..."
	@mkdir -p bin
	go build -o bin/om main.go
	@echo "✅ Build complete: bin/om"

# Build for all platforms
build-all:
	@echo "Building for all platforms..."
	@mkdir -p bin
	@if [ "$(OS)" = "Windows_NT" ]; then \
		powershell -ExecutionPolicy Bypass -File build-all.ps1; \
	else \
		chmod +x build-all.sh && ./build-all.sh; \
	fi

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	@echo "✅ Clean complete"

# Install dependencies
deps:
	@echo "Installing/updating dependencies..."
	go mod tidy
	go mod download
	@echo "✅ Dependencies updated"

# Run linters
lint:
	@echo "Running linters..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Format code
format:
	@echo "Formatting code..."
	go fmt ./...
	@echo "✅ Code formatted"

# Install binary to system PATH
install: build
	@echo "Installing binary..."
	@if [ "$(OS)" = "Windows_NT" ]; then \
		copy bin\om.exe C:\Windows\System32\om.exe; \
	else \
		sudo cp bin/om /usr/local/bin/om; \
	fi
	@echo "✅ Installation complete"

# Development setup
dev-setup: deps format lint test
	@echo "✅ Development setup complete"

# Quick development build
dev: deps build
	@echo "✅ Development build complete" 