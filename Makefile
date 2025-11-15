.PHONY: build test clean install release help

VERSION := $(shell cat VERSION)
BINARY_NAME := git-bits
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build development version
	mkdir -p bin
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) .

test: build ## Run test suite
	@export PATH="$(PWD)/bin:$(PATH)" && go test -v ./ ./bits ./command

clean: ## Clean build artifacts
	rm -rf bin/

install: ## Install dependencies
	go mod download
	go mod tidy

release: ## Cross compile release builds
	mkdir -p bin
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/linux_amd64/$(BINARY_NAME) .
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/windows_amd64/$(BINARY_NAME).exe .
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/darwin_amd64/$(BINARY_NAME) .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/darwin_arm64/$(BINARY_NAME) .


localstack-up: ## Start LocalStack for testing
	docker compose --file docker-compose.test.yml up --detach localstack
	@echo "Waiting for LocalStack to be ready..."
	@./scripts/setup-localstack.sh

localstack-down: ## Stop LocalStack
	docker compose --file docker-compose.test.yml down

docker-test: ## Run full test suite with LocalStack
	docker compose --file docker-compose.test.yml up --build --abort-on-container-exit --exit-code-from test-git-bits

docker-logs: ## Show LocalStack logs
	docker compose --file docker-compose.test.yml logs --follow

docker-clean: ## Clean up Docker test resources
	docker compose --file docker-compose.test.yml --verbose down
	docker builder prune --force
	#docker system prune --all --force --volumes