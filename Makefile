.PHONY: help install run migrate-create migrate-up migrate-status db-reset

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

install: ## Install dependencies
	go mod download
	go mod tidy

run: ## Run the application
	go run cmd/api/main.go

dev: ## Run with air (hot reload)
	air

migrate-create: ## Create a new migration (usage: make migrate-create name=add_books_table)
	atlas migrate diff $(name) --env local

migrate-up: ## Apply migrations
	atlas migrate apply --env local

migrate-status: ## Check migration status
	atlas migrate status --env local

migrate-hash: ## Hash/rehash migration directory
	atlas migrate hash --env local

db-reset: ## Reset database (WARNING: drops all data)
	dropdb mcs_dctfweb_sender || true
	createdb mcs_dctfweb_sender
	make migrate-up

test: ## Run tests
	go test -v ./...
