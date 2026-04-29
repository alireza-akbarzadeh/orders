# Makefile for orders-api

# Variables
BINARY_NAME=orders-api
DB_PATH=./orders.db
MIGRATIONS_DIR=./migrations
GOOSE_CMD=goose -dir $(MIGRATIONS_DIR) sqlite3 $(DB_PATH)

# Go related variables
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin

# Colors for output
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)

.PHONY: help build run clean test migrate create-migration reset-db deps tidy

# Default target
help: ## Show this help message
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  ${YELLOW}%-20s${RESET} ${GREEN}%s${RESET}\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

## Development
build: ## Build the application
	@echo "${GREEN}Building application...${RESET}"
	go build -o $(GOBIN)/$(BINARY_NAME) ./main.go
	@echo "${GREEN}Build complete: $(GOBIN)/$(BINARY_NAME)${RESET}"

run: ## Run the application
	@echo "${GREEN}Running application...${RESET}"
	go run ./main.go

clean: ## Clean build artifacts and database
	@echo "${YELLOW}Cleaning build artifacts...${RESET}"
	rm -rf $(GOBIN)
	@echo "${YELLOW}Cleaning database...${RESET}"
	rm -f $(DB_PATH)
	@echo "${GREEN}Clean complete${RESET}"

## Database Operations
create-migration: ## Create a new migration file (usage: make create-migration name=migration_name)
	@if [ -z "$(name)" ]; then \
		echo "${YELLOW}Error: Migration name required${RESET}"; \
		echo "${GREEN}Usage: make create-migration name=your_migration_name${RESET}"; \
		exit 1; \
	fi
	@echo "${GREEN}Creating migration: $(name)${RESET}"
	$(GOOSE_CMD) create $(name) sql

migrate-up: ## Run all pending migrations
	@echo "${GREEN}Running migrations...${RESET}"
	$(GOOSE_CMD) up
	@echo "${GREEN}Migrations completed${RESET}"

migrate-down: ## Rollback the last migration
	@echo "${YELLOW}Rolling back last migration...${RESET}"
	$(GOOSE_CMD) down
	@echo "${GREEN}Rollback completed${RESET}"

migrate-reset: ## Rollback all migrations
	@echo "${YELLOW}Resetting all migrations...${RESET}"
	$(GOOSE_CMD) reset
	@echo "${GREEN}Reset completed${RESET}"

migrate-status: ## Check migration status
	@echo "${GREEN}Migration status:${RESET}"
	$(GOOSE_CMD) status

migrate-version: ## Show current migration version
	@echo "${GREEN}Current migration version:${RESET}"
	$(GOOSE_CMD) version

reset-db: clean migrate-up ## Complete reset: clean + recreate database + run migrations
	@echo "${GREEN}Database reset complete${RESET}"

## Dependencies
deps: ## Install dependencies
	@echo "${GREEN}Installing dependencies...${RESET}"
	go mod download
	@echo "${GREEN}Installing goose...${RESET}"
	go install github.com/pressly/goose/v3/cmd/goose@latest
	@echo "${GREEN}Dependencies installed${RESET}"

tidy: ## Tidy up go.mod and go.sum
	@echo "${GREEN}Tidying go modules...${RESET}"
	go mod tidy
	@echo "${GREEN}Complete${RESET}"

## Testing
test: ## Run tests
	@echo "${GREEN}Running tests...${RESET}"
	go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "${GREEN}Running tests with coverage...${RESET}"
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "${GREEN}Coverage report generated: coverage.html${RESET}"

## Application Commands (Your specific workflow)
setup: deps create-orders-table run ## Complete setup: install deps, create orders table, run app

create-orders-table: ## Create orders_table migration and run it
	@echo "${GREEN}Creating orders_table migration...${RESET}"
	$(GOOSE_CMD) create orders_table sql
	@echo "${GREEN}Running migration...${RESET}"
	$(GOOSE_CMD) up
	@echo "${GREEN}Orders table created${RESET}"

## Utility
watch: ## Run app with hot reload (requires air)
	@command -v air >/dev/null 2>&1 || { \
		echo "${YELLOW}Installing air for hot reload...${RESET}"; \
		go install github.com/cosmtrek/air@latest; \
	}
	air

.PHONY: build run clean test migrate-up migrate-down migrate-reset migrate-status reset-db deps tidy help create-orders-table