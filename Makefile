.PHONY: build run test clean install docker-build docker-run fmt lint

BINARY_NAME=wireguard_exporter
VERSION?=1.0.0
BUILD_DIR=./build
INSTALL_DIR=/usr/local/bin

# Build variables
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

# Build the binary
build:
	@echo "🔨 Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/wireguard_exporter
	@echo "✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Run the exporter (requires sudo)
run: build
	@echo "🚀 Running $(BINARY_NAME)..."
	sudo $(BUILD_DIR)/$(BINARY_NAME)

# Run with verbose logging
run-verbose: build
	@echo "🚀 Running $(BINARY_NAME) (verbose mode)..."
	sudo $(BUILD_DIR)/$(BINARY_NAME) -verbose

# Run tests
test:
	@echo "🧪 Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	rm -rf $(BUILD_DIR)
	go clean
	@echo "✅ Clean complete"

# Install binary to system
install: build
	@echo "📦 Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/
	sudo chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "✅ Installation complete"

# Uninstall binary from system
uninstall:
	@echo "🗑️  Uninstalling $(BINARY_NAME)..."
	sudo rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "✅ Uninstallation complete"

# Format code
fmt:
	@echo "🎨 Formatting code..."
	go fmt ./...
	@echo "✅ Format complete"

# Run linter
lint:
	@echo "🔍 Running linter..."
	golangci-lint run
	@echo "✅ Lint complete"

# Build Docker image
docker-build:
	@echo "🐳 Building Docker image..."
	docker build -t $(BINARY_NAME):$(VERSION) -t $(BINARY_NAME):latest .
	@echo "✅ Docker build complete"

# Run Docker container
docker-run: docker-build
	@echo "🐳 Running Docker container..."
	docker run --rm --net=host --cap-add=NET_ADMIN $(BINARY_NAME):latest

# Build for multiple platforms
build-all:
	@echo "🔨 Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/wireguard_exporter
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/wireguard_exporter
	GOOS=linux GOARCH=arm go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm ./cmd/wireguard_exporter
	@echo "✅ Multi-platform build complete"

# Show help
help:
	@echo "Available targets:"
	@echo "  build         - Build the binary"
	@echo "  run           - Build and run the exporter"
	@echo "  run-verbose   - Build and run with verbose logging"
	@echo "  test          - Run tests"
	@echo "  clean         - Remove build artifacts"
	@echo "  install       - Install binary to $(INSTALL_DIR)"
	@echo "  uninstall     - Remove binary from $(INSTALL_DIR)"
	@echo "  fmt           - Format code"
	@echo "  lint          - Run linter"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Build and run Docker container"
	@echo "  build-all     - Build for multiple platforms"
	@echo "  help          - Show this help message"

