.PHONY: all help deps build shed install lint

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
	@echo "  make lint    - Format and lint code using golangci-lint"
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

lint:
	@echo "Formatting and linting code with golangci-lint..."
	@hash golangci-lint > /dev/null 2>&1 || (echo "Install golangci-lint to continue: https://golangci-lint.run/usage/install/"; exit 1)
	golangci-lint run --config .golangci.yaml
	@echo "✅ Code formatting and linting completed"

deps:
	@echo "Checking dependencies..."
	@hash go > /dev/null 2>&1 || (echo "Install go to continue: https://github.com/golang/go"; exit 1)
	@hash sqlc > /dev/null 2>&1 || (echo "Install sqlc to continue: https://github.com/sqlc-dev/sqlc"; exit 1)
	@hash docker > /dev/null 2>&1 || (echo "Install docker to continue: https://docs.docker.com/engine/install/"; exit 1)
	@hash migrate > /dev/null 2>&1 || (echo "Install golang-migrate to continue: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate"; exit 1)
	@echo "✅ All dependencies found"

sqlc: deps
	sqlc generate
