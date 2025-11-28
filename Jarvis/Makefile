.PHONY: all build test lint clean fmt help install

# Binary name
BINARY_NAME=jarvis

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOFMT=$(GOCMD) fmt
GOMOD=$(GOCMD) mod

all: test build

## build: Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) -o $(BINARY_NAME) -v

## test: Run all tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

## test-coverage: Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## test-race: Run tests with race detector
test-race:
	@echo "Running tests with race detector..."
	$(GOTEST) -race -v ./...

## bench: Run benchmarks
bench:
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

## lint: Run golangci-lint
lint:
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin"; \
	fi

## fmt: Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

## vet: Run go vet
vet:
	@echo "Running go vet..."
	$(GOCMD) vet ./...

## clean: Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

## verify: Verify dependencies
verify:
	@echo "Verifying dependencies..."
	$(GOMOD) verify

## install: Install the binary to GOPATH/bin
install: build
	@echo "Installing $(BINARY_NAME)..."
	cp $(BINARY_NAME) $(GOPATH)/bin/

## run: Build and run the binary
run: build
	./$(BINARY_NAME)

## check: Run all checks (fmt, vet, lint, test)
check: fmt vet lint test
	@echo "All checks passed!"

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'
