.PHONY: all help deps debug-build build shed install lint lint-fix lint-new test test-coverage

COMMIT := $(shell git rev-parse --short HEAD)
GOOS := $(shell go env GOOS)
GOARCH := arm64
BINARY_NAME := shed

# Detect OS and set build tag (linux, darwin, or windows)
UNAME_S := $(shell uname -s 2>/dev/null || echo "unknown")
ifeq ($(UNAME_S),Linux)
	OS_TAG := linux
else ifeq ($(UNAME_S),Darwin)
	OS_TAG := darwin
else ifeq ($(findstring MINGW,$(UNAME_S)),MINGW)
	OS_TAG := windows
else ifeq ($(findstring MSYS,$(UNAME_S)),MSYS)
	OS_TAG := windows
else ifeq ($(findstring CYGWIN,$(UNAME_S)),CYGWIN)
	OS_TAG := windows
else
	OS_TAG := $(GOOS)
endif

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
	go build -tags="sqlcipher,$(OS_TAG)" -ldflags="-X h3jfc/shed/cmd.Commit=$(COMMIT)" -o ./$(BINARY_NAME) main.go
	@echo "✅ Built $(BINARY_NAME) successfully"

debug-build: deps
	@echo "Building a debug build of $(BINARY_NAME) with SQLCipher support..."
	go build -gcflags=all="-N -l" -tags="sqlcipher,$(OS_TAG)" -ldflags="-X h3jfc/shed/cmd.Commit=$(COMMIT)" -o ./$(BINARY_NAME) main.go
	@echo "✅ Built $(BINARY_NAME) successfully"

shed: deps
	@echo "Running application with SQLCipher support..."
	go run -tags="sqlcipher,$(OS_TAG)" -ldflags="-X h3jfc/shed/main.Commit=$(COMMIT)" main.go

install: deps
	@echo "Building and installing $(BINARY_NAME) with SQLCipher support to system..."
	go install -tags="sqlcipher,$(OS_TAG)" -ldflags="-X h3jfc/shed/main.Commit=$(COMMIT)" .
	@echo "✅ Installed $(BINARY_NAME) to system PATH"

lint: deps
	@echo "Formatting and linting code with golangci-lint..."
	golangci-lint run --config .golangci.yaml --build-tags=sqlcipher,$(OS_TAG)
	wait
	@echo "✅ Code formatting and linting completed"

lint-fix: deps
	@echo "Formatting and linting code with golangci-lint (auto-fix enabled, parallel)..."
	golangci-lint run --fix --config .golangci.yaml --build-tags=sqlcipher,$(OS_TAG) --allow-parallel-runners &
	wait
	@echo "✅ Code formatting and linting completed with auto-fix"

lint-parallel: deps
	@echo "Formatting and linting code with golangci-lint (parallel)..."
	golangci-lint run --config .golangci.yaml --allow-parallel-runners &
	wait
	@echo "✅ Code formatting and linting completed"

lint-new: deps
	@echo "Linting only new/changed code with golangci-lint..."
	golangci-lint run --new --config .golangci.yaml --allow-parallel-runners
	@echo "✅ New code linting completed"

test: deps
	@echo "Running all Go tests..."
	go test -tags="sqlcipher,$(OS_TAG)" -ldflags="-X h3jfc/shed/main.Commit=$(COMMIT)" -v ./...
	@echo "✅ All tests completed"

test-coverage: deps
	@echo "Running tests with coverage analysis..."
	go test -tags="sqlcipher,$(OS_TAG)" -coverprofile=coverage.out -covermode=atomic ./...
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
