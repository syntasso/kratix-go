# kratix-go Makefile
# Consistent with syntasso/kratix repository structure

.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: tidy
tidy: ## Clean up go.mod and go.sum
	go mod tidy
	go mod verify

.PHONY: fmt
fmt: ## Format Go code
	go fmt ./...

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: lint
lint: fmt vet ## Run all linting checks

.PHONY: build
build: ## Build the Go package
	go build .

##@ Testing

.PHONY: test
test: ## Run unit tests
	go test -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: e2e-test
e2e-test: build-and-load-image ## Run end-to-end tests using Ginkgo
	ginkgo -v -r system

.PHONY: ginkgo-test
ginkgo-test: ## Run all tests using Ginkgo
	ginkgo -v

.PHONY: test-watch
test-watch: ## Run tests in watch mode
	ginkgo watch -v

##@ CI/CD

.PHONY: ci-test
ci-test: lint test ## Run all CI tests (lint + unit tests)

.PHONY: clean
clean: ## Clean build artifacts and test files
	go clean
	rm -f coverage.out coverage.html

##@ Utilities

.PHONY: deps
deps: ## Download and verify dependencies
	go mod download
	go mod verify

.PHONY: upgrade
upgrade: ## Upgrade all dependencies
	go get -u ./...
	go mod tidy

build-and-load-image:
	docker build -t ghcr.io/syntasso/kratix-go/sdk-test:v1.0.0 -f system/assets/workflows/Dockerfile .
	kind load docker-image ghcr.io/syntasso/kratix-go/sdk-test:v1.0.0 --name platform
