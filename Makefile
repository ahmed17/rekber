# File: Makefile
.PHONY: help setup run build test clean migrate-up migrate-down dev db-create db-drop deps tidy watch seed test-db

# Colors for output
GREEN  := \033[0;32m
YELLOW := \033[1;33m
RED    := \033[0;31m
NC     := \033[0m

# Default target
.DEFAULT_GOAL := help

## Help: Menampilkan semua command yang tersedia
help:
	@echo "${YELLOW}Available commands:${NC}"
	@echo ""
	@echo "${YELLOW}Development:${NC}"
	@echo "  ${GREEN}make setup${NC}     - Setup project (install dependencies and tools)"
	@echo "  ${GREEN}make dev${NC}       - Run development server with hot reload"
	@echo "  ${GREEN}make run${NC}       - Run application normally"
	@echo ""
	@echo "${YELLOW}Database:${NC}"
	@echo "  ${GREEN}make db-create${NC}    - Create database (port 5433)"
	@echo "  ${GREEN}make db-drop${NC}      - Drop database"
	@echo "  ${GREEN}make db-reset${NC}     - Drop and recreate database"
	@echo "  ${GREEN}make migrate-up${NC}   - Run all migrations"
	@echo "  ${GREEN}make migrate-down${NC} - Rollback last migration"
	@echo "  ${GREEN}make migrate-new${NC}  - Create new migration file"
	@echo "  ${GREEN}make test-db${NC}      - Test database connection"
	@echo ""
	@echo "${YELLOW}Testing:${NC}"
	@echo "  ${GREEN}make test${NC}       - Run all tests"
	@echo "  ${GREEN}make test-cover${NC} - Run tests with coverage"
	@echo "  ${GREEN}make lint${NC}       - Run linter"
	@echo ""
	@echo "${YELLOW}Maintenance:${NC}"
	@echo "  ${GREEN}make clean${NC}      - Clean build artifacts"
	@echo "  ${GREEN}make deps${NC}       - Download all dependencies"
	@echo "  ${GREEN}make tidy${NC}       - Tidy go.mod"
	@echo ""

## Setup: Install dependencies and tools
setup:
	@echo "${GREEN}Installing dependencies...${NC}"
	go mod download
	@echo "${GREEN}Installing air for hot reload...${NC}"
	go install github.com/cosmtrek/air@latest
	@echo "${GREEN}Installing migrate tool...${NC}"
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "${GREEN}Setup completed!${NC}"

## Run: Run the application
run:
	@echo "${GREEN}Starting server...${NC}"
	go run cmd/api/main.go

## Dev: Run with hot reload (using air)
dev:
	@if command -v air > /dev/null; then \
		echo "${GREEN}Starting with air (hot reload)...${NC}"; \
		air; \
	else \
		echo "${YELLOW}Air not found, installing...${NC}"; \
		go install github.com/cosmtrek/air@latest; \
		air; \
	fi

## Database: Create database
db-create:
	@echo "${GREEN}Creating database 'db_rekber' on port 5433...${NC}"
	PGPASSWORD=password createdb -h localhost -p 5433 -U postgres db_rekber 2>/dev/null || echo "${YELLOW}Database already exists or error occurred${NC}"
	@echo "${GREEN}Database 'db_rekber' ready${NC}"

## Database: Drop database
db-drop:
	@echo "${YELLOW}Dropping database 'db_rekber'...${NC}"
	PGPASSWORD=password dropdb -h localhost -p 5433 -U postgres db_rekber 2>/dev/null || echo "${YELLOW}Database doesn't exist or error occurred${NC}"
	@echo "${GREEN}Database dropped${NC}"

## Database: Reset (drop and recreate)
db-reset: db-drop db-create migrate-up

## Test database connection
test-db:
	@echo "${GREEN}Testing database connection...${NC}"
	go run scripts/test-connection.go

## Migrate: Run migrations (FIXED SSL)
migrate-up:
	@echo "${GREEN}Running migrations...${NC}"
	migrate -path migrations -database "postgres://postgres:password@localhost:5433/db_rekber?sslmode=disable" up

## Migrate: Rollback last migration
migrate-down:
	@echo "${GREEN}Rolling back migration...${NC}"
	migrate -path migrations -database "postgres://postgres:password@localhost:5433/db_rekber?sslmode=disable" down 1

## Migrate: Create new migration file
migrate-new:
	@echo "${GREEN}Creating new migration...${NC}"
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $${name// /_}

## Seed: Seed database with test data
seed:
	@echo "${GREEN}Seeding database...${NC}"
	@if [ -f scripts/seed.go ]; then \
		go run scripts/seed.go; \
		echo "${GREEN}Database seeded successfully${NC}"; \
	else \
		echo "${RED}Seed file not found${NC}"; \
	fi

## Clean: Remove build artifacts
clean:
	@echo "${GREEN}Cleaning...${NC}"
	rm -rf bin/
	rm -rf coverage.out coverage.html
	rm -rf uploads/
	go clean

## Dependencies: Download all dependencies
deps:
	@echo "${GREEN}Downloading dependencies...${NC}"
	go mod download

## Tidy: Clean up go.mod
tidy:
	@echo "${GREEN}Tidying go.mod...${NC}"
	go mod tidy

## Test: Run all tests
test:
	@echo "${GREEN}Running tests...${NC}"
	go test ./... -v

## Quick start for development
quick-start:
	@echo "${GREEN}=== Quick Start for Development ===${NC}"
	@echo "1. Creating database..."
	@$(MAKE) db-create
	@echo "2. Running migrations..."
	@$(MAKE) migrate-up
	@echo "3. Starting development server..."
	@$(MAKE) dev