.PHONY: help clean build run lint test fmt vet update-deps

# Variables
BINARY_NAME=eval-hub-backend-svc
CMD_PATH=./cmd/eval_hub
BIN_DIR=bin
PORT?=8080

# Default target
.DEFAULT_GOAL := help

help: ## Display this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

clean: ## Remove build artifacts
	@echo "Cleaning..."
	@rm -rf $(BIN_DIR)
	@rm -f $(BINARY_NAME)
	@rm -f coverage.out coverage.html
	@go clean
	@echo "Clean complete"

build: ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/$(BINARY_NAME) $(CMD_PATH)
	@echo "Build complete: $(BIN_DIR)/$(BINARY_NAME)"

run: ## Run the application
	@echo "Running $(BINARY_NAME) on port $(PORT)..."
	@PORT=$(PORT) go run $(CMD_PATH)/main.go

lint: ## Lint the code (runs go vet, skips go fmt to preserve 2-space indentation)
	@echo "Linting code..."
	@go vet ./...
	@echo "Lint complete"

fmt: ## Format the code with go fmt (NOTE: converts indentation to tabs per Go standard)
	@echo "Formatting code with go fmt..."
	@echo "WARNING: go fmt will convert indentation to tabs (Go standard)"
	@go fmt ./...
	@echo "Format complete"

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...
	@echo "Vet complete"

test: ## Run unit tests
	@echo "Running unit tests..."
	@go test -v ./internal/...

test-fvt: ## Run FVT (Functional Verification Tests) using godog
	@echo "Running FVT tests..."
	@go test -v ./tests/features/...

test-all: test test-fvt ## Run all tests (unit + FVT)

test-coverage: ## Run unit tests with coverage
	@echo "Running unit tests with coverage..."
	@go test -v -coverprofile=coverage.out ./internal/...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

install-deps: ## Install dependencies
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies installed"

update-deps: ## Update all dependencies to latest versions
	@echo "Updating dependencies to latest versions..."
	@go get -u ./...
	@go mod tidy
	@echo "Dependencies updated"
