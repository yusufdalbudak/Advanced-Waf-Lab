.PHONY: build test lint run clean docker docker-run help

# Variables
BINARY_NAME=wafd
DOCKER_IMAGE=waf:latest
GO_VERSION=1.21

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the WAF binary
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) ./cmd/wafd
	@echo "Build complete!"

test: ## Run tests
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

test-coverage: ## Run tests with coverage
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out

lint: ## Run linters
	@echo "Running linters..."
	@golangci-lint run

lint-fix: ## Run linters and fix issues
	@golangci-lint run --fix

fmt: ## Format code
	@go fmt ./...
	@goimports -w .

run: build ## Build and run WAF
	@./$(BINARY_NAME)

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -f test-website dashboard
	@rm -f coverage.out coverage.html
	@rm -f *.log
	@echo "Clean complete!"

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE) .
	@echo "Docker image built: $(DOCKER_IMAGE)"

docker-run: docker-build ## Build and run Docker container
	@echo "Running Docker container..."
	@docker run -d \
		--name waf \
		-p 8080:8080 \
		-v $(PWD)/configs:/app/configs \
		-v $(PWD)/logs:/app/logs \
		$(DOCKER_IMAGE)
	@echo "Container running on http://localhost:8080"

docker-stop: ## Stop Docker container
	@docker stop waf || true
	@docker rm waf || true

docker-logs: ## View Docker container logs
	@docker logs -f waf

benchmark: ## Run benchmarks
	@go test -bench=. -benchmem ./...

fuzz: ## Run fuzz tests
	@go test -fuzz=. ./test/fuzz/...

security-scan: ## Run security scan
	@gosec ./...

install-tools: ## Install development tools
	@echo "Installing tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install github.com/securego/gosec/v2/cmd/gosec@latest
	@echo "Tools installed!"

ci: lint test ## Run CI checks (lint + test)

all: clean lint test build ## Run all checks and build

