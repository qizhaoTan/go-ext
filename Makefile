.PHONY: help test cover lint fmt vet check tidy

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2}'

test: ## Run the tests with the race detector
	go test -race ./...

cover: ## Run the tests and open the coverage report
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...
	go tool cover -html=coverage.txt -o coverage.html
	@echo "coverage report written to coverage.html"

lint: ## Run golangci-lint
	golangci-lint run

fmt: ## Format the code
	gofmt -w .

vet: ## Run go vet
	go vet ./...

check: fmt vet test ## Format, vet and test

tidy: ## Tidy go.mod
	go mod tidy
