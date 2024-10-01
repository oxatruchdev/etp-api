# Makefile

MIGRATION_DIR = db/migrations
DATE := $(shell date +%Y%m%d%H%M%S)

build:
	@echo "Building..."
	@templ generate
	@npx tailwindcss -i cmd/web/assets/css/input.css -o cmd/web/assets/css/output.css
	@go build -o main cmd/etp/main.go

# Create a new migration file with a timestamp-based name.
migration:
	@read -p "Enter migration description: " desc; \
	name=$(DATE)_$${desc// /_}.sql; \
	touch $(MIGRATION_DIR)/$$name; \
	echo "Created migration file: $(MIGRATION_DIR)/$$name"
