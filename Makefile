.PHONY: help install run dev test \
        grove-build grove-install \
        migrate-create migrate-up migrate-status migrate-hash db-reset

# ──────────────────────────────────────────────────────────────────────────────
# Help
# ──────────────────────────────────────────────────────────────────────────────

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-22s\033[0m %s\n", $$1, $$2}'

# ──────────────────────────────────────────────────────────────────────────────
# Dependencies
# ──────────────────────────────────────────────────────────────────────────────

install: ## Install Go dependencies
	go mod download
	go mod tidy

# ──────────────────────────────────────────────────────────────────────────────
# Grove CLI
# ──────────────────────────────────────────────────────────────────────────────

grove-build: ## Build the grove CLI binary → bin/grove
	@mkdir -p bin
	go build -o bin/grove ./cmd/grove
	@echo "\033[32m  ✔ Binary built at bin/grove\033[0m"

grove-install: ## Install grove globally (go install → adds to PATH)
	go install ./cmd/grove
	@echo "\033[32m  ✔ grove installed — run: grove --help\033[0m"

# ──────────────────────────────────────────────────────────────────────────────
# Application
# ──────────────────────────────────────────────────────────────────────────────

run: ## Run the application (go run)
	go run ./cmd/api/main.go

dev: ## Run with air (hot reload)
	air

serve: grove-build ## Build grove then start the dev server (grove serve)
	./bin/grove serve

test: ## Run all tests
	go test -v ./...

# ──────────────────────────────────────────────────────────────────────────────
# Migrations (legacy make targets — prefer: grove migrate)
# ──────────────────────────────────────────────────────────────────────────────

migrate-create: ## Create a migration  (usage: make migrate-create name=add_posts_table)
	atlas migrate diff $(name) --env local

migrate-up: ## Apply all pending migrations
	atlas migrate apply --env local

migrate-status: ## Show migration status
	atlas migrate status --env local

migrate-hash: ## Rehash the migrations directory
	atlas migrate hash --env local

db-reset: ## Reset DB — drops all tables and re-applies migrations (DEV ONLY)
	atlas schema clean --env local --auto-approve
	atlas migrate apply --env local
