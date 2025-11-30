.PHONY: help test test-verbose test-coverage test-coverage-html run build clean lint fmt vet tidy install-deps

# Variables
BINARY_NAME=recipes-api
COVERAGE_FILE=coverage.out
GO=go
GOTEST=$(GO) test
GOBUILD=$(GO) build
GOCLEAN=$(GO) clean
GOGET=$(GO) get
GOMOD=$(GO) mod
GOFMT=$(GO) fmt
GOVET=$(GO) vet

# Default target
help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

test: ## Run all tests
	@echo "Running tests..."
	@$(GOTEST) ./... -v

test-short: ## Run tests without verbose output
	@echo "Running tests (short)..."
	@$(GOTEST) ./...

test-race: ## Run tests with race detector
	@echo "Running tests with race detector..."
	@$(GOTEST) ./... -race

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	@$(GOTEST) ./... -coverprofile=$(COVERAGE_FILE)
	@$(GO) tool cover -func=$(COVERAGE_FILE)
	@echo ""
	@echo "Coverage report saved to $(COVERAGE_FILE)"

test-coverage-html: test-coverage ## Generate HTML coverage report
	@echo "Generating HTML coverage report..."
	@$(GO) tool cover -html=$(COVERAGE_FILE)

test-watch: ## Run tests continuously (requires entr)
	@echo "Watching for changes..."
	@find . -name '*.go' | entr -c make test-short

run: ## Run the application
	@echo "Starting application..."
	@$(GO) run main.go

build: ## Build the application
	@echo "Building $(BINARY_NAME)..."
	@$(GOBUILD) -o $(BINARY_NAME) -v .
	@echo "Binary created: $(BINARY_NAME)"

build-linux: ## Build for Linux
	@echo "Building for Linux..."
	@GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)-linux-amd64 -v .

build-windows: ## Build for Windows
	@echo "Building for Windows..."
	@GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)-windows-amd64.exe -v .

build-mac: ## Build for macOS
	@echo "Building for macOS..."
	@GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)-darwin-amd64 -v .

build-all: build-linux build-windows build-mac ## Build for all platforms

clean: ## Clean build artifacts and coverage files
	@echo "Cleaning..."
	@$(GOCLEAN)
	@rm -f $(BINARY_NAME)
	@rm -f $(BINARY_NAME)-*
	@rm -f $(COVERAGE_FILE)
	@echo "Clean complete"

lint: ## Run golangci-lint (requires golangci-lint to be installed)
	@echo "Running linter..."
	@golangci-lint run ./...

fmt: ## Format code
	@echo "Formatting code..."
	@$(GOFMT) ./...

vet: ## Run go vet
	@echo "Running go vet..."
	@$(GOVET) ./...

tidy: ## Tidy go modules
	@echo "Tidying modules..."
	@$(GOMOD) tidy

install-deps: ## Install/update dependencies
	@echo "Installing dependencies..."
	@$(GOMOD) download
	@$(GOMOD) verify

check: fmt vet test ## Run format, vet, and tests

ci: check test-coverage ## Run all CI checks

all: clean install-deps check build ## Run clean, install deps, checks, and build

# Development helpers
dev: ## Run in development mode with hot reload (requires air)
	@echo "Starting development server with hot reload..."
	@air

install-tools: ## Install development tools
	@echo "Installing development tools..."
	@$(GO) install github.com/cosmtrek/air@latest
	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Tools installed successfully"

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t $(BINARY_NAME):latest .

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	@docker run -p 8080:8080 $(BINARY_NAME):latest

.DEFAULT_GOAL := help
