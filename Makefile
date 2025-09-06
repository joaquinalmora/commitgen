# commitgen Makefile
# Common development tasks

.PHONY: build test install clean doctor release dev-setup help

# Default target
all: build

# Build the binary
build:
	@echo "Building commitgen..."
	go build -o bin/commitgen ./cmd/commitgen
	@echo "✅ Built bin/commitgen"

# Run all tests
test:
	@echo "Running tests..."
	go test ./...
	@echo "✅ Tests completed"

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

# Install to system
install:
	@echo "Installing commitgen..."
	go install ./cmd/commitgen
	@echo "✅ Installed commitgen to $(shell go env GOPATH)/bin"

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	@echo "✅ Cleaned up build artifacts"

# Run system health check
doctor:
	@echo "Running system diagnostics..."
	./bin/commitgen doctor

# Development setup
dev-setup:
	@echo "Setting up development environment..."
	go mod tidy
	go mod download
	@echo "✅ Development environment ready"

# Quick development build and test
dev: build test
	@echo "✅ Development build and test completed"

# Show available targets
help:
	@echo "Available targets:"
	@echo "  build         - Build the binary"
	@echo "  test          - Run all tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  install       - Install to system"
	@echo "  clean         - Clean build artifacts"
	@echo "  doctor        - Run system health check"
	@echo "  dev-setup     - Setup development environment"
	@echo "  dev           - Quick build and test"
	@echo "  release-dry   - Test release process (dry run)"
	@echo "  help          - Show this help"

# Test release process (dry run)
release-dry:
	@echo "Testing release process..."
	@command -v goreleaser >/dev/null 2>&1 || { echo "Installing GoReleaser..."; curl -sfL https://goreleaser.com/static/run | bash -s -- --version; }
	goreleaser release --snapshot --clean
	@echo "✅ Release dry run completed - check dist/ directory"
