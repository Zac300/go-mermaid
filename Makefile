.PHONY: build test cover lint fmt golden tidy clean

build: ## Build the CLI
	go build -o bin/mermaid ./cmd/mermaid

test: ## Run tests with race detector
	go test -race ./...

cover: ## Run tests and open an HTML coverage report
	go test -coverprofile=coverage.out -coverpkg=./... ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "wrote coverage.html"

lint: ## Run golangci-lint (must be installed)
	golangci-lint run

fmt: ## Format code
	gofmt -s -w .
	go run golang.org/x/tools/cmd/goimports@latest -w .

golden: ## Regenerate golden SVG files
	go test -run TestGolden -update .

tidy: ## Tidy modules
	go mod tidy

clean: ## Remove build artifacts
	rm -rf bin coverage.out coverage.html

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'
