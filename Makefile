# Makefile for pdfunlock

# Binary name
BINARY_NAME=pdfunlock

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags "-s -w"

# Platform detection
UNAME_S := $(shell uname -s)
UNAME_M := $(shell uname -m)
ifeq ($(UNAME_S),Linux)
    PLATFORM=linux
endif
ifeq ($(UNAME_S),Darwin)
    PLATFORM=darwin
endif
ifeq ($(UNAME_S),MINGW32_NT)
    PLATFORM=windows
endif
ifeq ($(UNAME_S),MINGW64_NT)
    PLATFORM=windows
endif
ifeq ($(UNAME_M),x86_64)
    ARCH=amd64
endif
ifeq ($(UNAME_M),arm64)
    ARCH=arm64
endif
ifeq ($(UNAME_M),aarch64)
    ARCH=arm64
endif

# Build directory
BUILD_DIR=bin/$(PLATFORM)_$(ARCH)
BUILD_PATH=$(BUILD_DIR)/$(BINARY_NAME)

.PHONY: all build clean test deps install uninstall build-all tag release help

# Default target
all: build

# Build the binary
build:
	@echo "Building $(BINARY_NAME) for $(PLATFORM)_$(ARCH)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_PATH) .
	@echo "Build complete: $(BUILD_PATH)"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf bin/
	@echo "Clean complete"

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOGET) -d ./...
	$(GOMOD) tidy

# Install binary to system
install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin/"
	@sudo cp $(BUILD_PATH) /usr/local/bin/
	@echo "Installation complete"

# Uninstall binary from system
uninstall:
	@echo "Removing $(BINARY_NAME) from /usr/local/bin/"
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "Uninstall complete"

# Build for multiple platforms using GoReleaser
build-all:
	@echo "Building for multiple platforms using GoReleaser..."
	@if ! command -v goreleaser >/dev/null 2>&1; then \
		echo "GoReleaser not found. Installing..."; \
		go install github.com/goreleaser/goreleaser@latest; \
	fi
	@if ! command -v upx >/dev/null 2>&1; then \
		echo "UPX not found. Please install UPX for compression support."; \
	fi
	@if [ -f .env ]; then \
		set -a && . ./.env && set +a && \
		goreleaser build --snapshot --clean; \
	else \
		GITHUB_REPOSITORY_OWNER=fadeltd \
		ENABLE_UPX_LINUX=true ENABLE_UPX_WINDOWS=true ENABLE_UPX_DARWIN_INTEL=true ENABLE_UPX_DARWIN_ARM=false \
		goreleaser build --snapshot --clean; \
	fi
	@echo "Multi-platform build complete with GoReleaser in bin/ directory"

# Create and push a new tag for release
# Usage: make tag VERSION=v1.0.0
tag:
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: VERSION is required. Usage: make tag VERSION=v1.0.0"; \
		exit 1; \
	fi
	@echo "Creating tag $(VERSION)..."
	@git tag -a $(VERSION) -m "Release $(VERSION)"
	@echo "Tag $(VERSION) created successfully"
	@echo "To push the tag and trigger release, run: make release VERSION=$(VERSION)"

# Push tag to trigger GitHub Actions release
# Usage: make release VERSION=v1.0.0
release:
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: VERSION is required. Usage: make release VERSION=v1.0.0"; \
		exit 1; \
	fi
	@echo "Pushing tag $(VERSION) to trigger release..."
	@git push origin $(VERSION)
	@echo "Release triggered! Check GitHub Actions for build status."
	@echo "Release will be available at: https://github.com/fadeltd/pdfunlock/releases/tag/$(VERSION)"

# Create and push tag in one command
# Usage: make tag-release VERSION=v1.0.0
tag-release: tag release

# Show help
help:
	@echo "Available targets:"
	@echo "  build        - Build the binary"
	@echo "  clean        - Clean build artifacts"
	@echo "  test         - Run tests"
	@echo "  deps         - Download dependencies"
	@echo "  install      - Install binary to system"
	@echo "  uninstall    - Remove binary from system"
	@echo "  build-all    - Build for multiple platforms"
	@echo "  tag          - Create a new tag (requires VERSION=vX.Y.Z)"
	@echo "  release      - Push tag to trigger release (requires VERSION=vX.Y.Z)"
	@echo "  tag-release  - Create and push tag in one command (requires VERSION=vX.Y.Z)"
	@echo "  help         - Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make tag VERSION=v1.0.0"
	@echo "  make release VERSION=v1.0.0"
	@echo "  make tag-release VERSION=v1.0.0"