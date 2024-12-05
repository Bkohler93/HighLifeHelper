DB_DIR = data/databases

init:
	@echo "Downloading dependencies and configuring project"
	@go mod download && go mod tidy
	@if [ ! -f "$(DB_DIR)/development/session/session.db" ]; then \
		mkdir -p "$(DB_DIR)/development/session"; \
		touch "$(DB_DIR)/development/session/session.db"; \
	fi
	@if [ ! -f "$(DB_DIR)/production/session/session.db" ]; then \
		mkdir -p "$(DB_DIR)/production/session"; \
		touch "$(DB_DIR)/production/session/session.db"; \
	fi
	@if [ ! -f .env ]; then \
		echo "No .env file found! Please create one based on .env.example."; \
		exit 1; \
	fi
	@echo ".env file found. Proceeding with setup..."
	@echo "Starting PostgreSQL container..."
	@docker compose up -d db
	@echo "Project set up!"

css:
	@echo "Minifying css"
	@./tailwindcss -i static/styles/input.css -o static/styles/output.css --minify

css-watch:
	@echo "Watching HTML files for changes"
	@./tailwindcss -i static/styles/input.css -o static/styles/output.css --watch

run:
	@go run .

clean:
	@docker compose down

migrate-down:
	@go run data/sql/migrations/migrateDown.go