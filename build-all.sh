#!/bin/bash

# Build script for om binaries - Cross-platform version
echo "Building om binaries for all platforms..."

# Create bin directory if it doesn't exist
mkdir -p bin

# For macOS (Intel)
echo "Building for macOS (Intel)..."
GOOS=darwin GOARCH=amd64 go build -o bin/om-darwin-amd64 main.go
if [ $? -eq 0 ]; then
    echo "✅ macOS (Intel) build successful"
else
    echo "❌ macOS (Intel) build failed"
fi

# For macOS (Apple Silicon)
echo "Building for macOS (Apple Silicon)..."
GOOS=darwin GOARCH=arm64 go build -o bin/om-darwin-arm64 main.go
if [ $? -eq 0 ]; then
    echo "✅ macOS (Apple Silicon) build successful"
else
    echo "❌ macOS (Apple Silicon) build failed"
fi

# For Linux
echo "Building for Linux..."
GOOS=linux GOARCH=amd64 go build -o bin/om-linux-amd64 main.go
if [ $? -eq 0 ]; then
    echo "✅ Linux build successful"
else
    echo "❌ Linux build failed"
fi

# For Windows (AMD64)
echo "Building for Windows (AMD64)..."
GOOS=windows GOARCH=amd64 go build -o bin/om-windows-amd64.exe main.go
if [ $? -eq 0 ]; then
    echo "✅ Windows (AMD64) build successful"
else
    echo "❌ Windows (AMD64) build failed"
fi

# For Windows (ARM64)
echo "Building for Windows (ARM64)..."
GOOS=windows GOARCH=arm64 go build -o bin/om-windows-arm64.exe main.go
if [ $? -eq 0 ]; then
    echo "✅ Windows (ARM64) build successful"
else
    echo "❌ Windows (ARM64) build failed"
fi

echo "Build process completed!"
echo "Binaries are available in the bin/ directory" 