.PHONY: all help deps

COMMIT := $(shell git rev-parse --short HEAD)
GOOS := linux
GOARCH := arm64

all: help

help:
	@echo "Available targets:"
	@echo "  make deps"

deps:
	@echo "Checking dependencies..."
	@hash go > /dev/null 2>&1 || (echo "Install go to continue: https://github.com/golang/go"; exit 1)
	@hash sqlc > /dev/null 2>&1 || (echo "Install sqlc to continue: https://github.com/sqlc-dev/sqlc"; exit 1)
	@hash docker > /dev/null 2>&1 || (echo "Install docker to continue: https://docs.docker.com/engine/install/"; exit 1)
	@hash migrate > /dev/null 2>&1 || (echo "Install golang-migrate to continue: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate"; exit 1)
	@echo "✅ All dependencies found"

sqlc: deps
	sqlc generate
