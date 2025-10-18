.PHONY: all help deps build shed install lint lint-fix lint-new test test-coverage

COMMIT := $(shell git rev-parse --short HEAD)
GOOS := darwin
GOARCH := arm64
BINARY_NAME := shed

all: install

help:
	@echo "Available targets:"
	@echo "  make         - Build and install shed to system (same as make install)"
	@echo "  make install - Build and install shed to system PATH"
	@echo "  make deps    - Check dependencies"
	@echo "  make build   - Build the application with SQLCipher support to ./shed"
	@echo "  make shed    - Run the application with SQLCipher support"
	@echo "  make test    - Run all Go tests"
	@echo "  make test-coverage - Run tests with coverage analysis"
	@echo "  make lint    - Format and lint code using golangci-lint (parallel)"
	@echo "  make lint-fix - Format and lint code with auto-fix enabled (parallel)"
	@echo "  make lint-new - Lint only new/changed code"
	@echo "  make sqlc    - Generate sqlc code"

build: deps
	@echo "Building $(BINARY_NAME) with SQLCipher support..."
	go build -tags="sqlcipher" -o ./$(BINARY_NAME) main.go
	@echo "✅ Built $(BINARY_NAME) successfully"

shed: deps
	@echo "Running application with SQLCipher support..."
	go run -tags="sqlcipher" main.go

install: deps
	@echo "Building and installing $(BINARY_NAME) with SQLCipher support to system..."
	go install -tags="sqlcipher" .
	@echo "✅ Installed $(BINARY_NAME) to system PATH"

lint: deps
	@echo "Formatting and linting code with golangci-lint (parallel)..."
	golangci-lint run --config .golangci.yaml --allow-parallel-runners &
	wait
	@echo "✅ Code formatting and linting completed"

lint-fix: deps
	@echo "Formatting and linting code with golangci-lint (auto-fix enabled, parallel)..."
	golangci-lint run --fix --config .golangci.yaml --allow-parallel-runners &
	wait
	@echo "✅ Code formatting and linting completed with auto-fix"

lint-new: deps
	@echo "Linting only new/changed code with golangci-lint..."
	golangci-lint run --new --config .golangci.yaml --allow-parallel-runners
	@echo "✅ New code linting completed"

test: deps
	@echo "Running all Go tests..."
	go test -tags="sqlcipher" -v ./...
	@echo "✅ All tests completed"

test-coverage: deps
	@echo "Running tests with coverage analysis..."
	go test -tags="sqlcipher" -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage analysis completed"
	@echo "📊 Coverage report: coverage.html"
	@echo "📋 Coverage profile: coverage.out"

deps:
	@echo "Checking dependencies..."
	@hash go > /dev/null 2>&1 || (echo "Install go to continue: https://github.com/golang/go"; exit 1)
	@hash sqlc > /dev/null 2>&1 || (echo "Install sqlc to continue: https://github.com/sqlc-dev/sqlc"; exit 1)
	@hash docker > /dev/null 2>&1 || (echo "Install docker to continue: https://docs.docker.com/engine/install/"; exit 1)
	@hash migrate > /dev/null 2>&1 || (echo "Install golang-migrate to continue: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate"; exit 1)
	@hash golangci-lint > /dev/null 2>&1 || (echo "Install golangci-lint to continue: https://golangci-lint.run/usage/install/"; exit 1)
	@echo "✅ All dependencies found"

sqlc: deps
	sqlc generate
