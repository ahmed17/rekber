.PHONY: help run build migrate-up migrate-down test clean

help:
	@echo "Available commands:"
	@echo "  make run       - Run the application in development mode"
	@echo "  make build     - Build the application"
	@echo "  make migrate-up   - Run database migrations"
	@echo "  make migrate-down - Rollback database migrations"
	@echo "  make test      - Run tests"
	@echo "  make clean     - Clean build artifacts"

run:
	@go run cmd/api/main.go

build:
	@go build -o bin/rekber cmd/api/main.go

migrate-up:
	@echo "Running migrations..."
	# Will be implemented later

migrate-down:
	@echo "Rolling back migrations..."
	# Will be implemented later

test:
	@go test ./...

clean:
	@rm -rf bin/ coverage.out

dev:
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Installing air for live reload..."; \
		go install github.com/cosmtrek/air@latest; \
		air; \
	fi