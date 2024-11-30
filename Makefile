DB_DIR = data/databases

init:
	@echo "Downloading dependencies and configuring project"
	go mod download && go mod tidy
	@if [ ! -f "$(DB_DIR)/development/session/session.db" ]; then \
		mkdir -p "$(DB_DIR)/development/session"; \
		touch "$(DB_DIR)/development/session/session.db"; \
	fi
	@if [ ! -f "$(DB_DIR)/production/session/session.db" ]; then \
		mkdir -p "$(DB_DIR)/production/session"; \
		touch "$(DB_DIR)/production/session/session.db"; \
	fi
	init:
	@echo "Downloading dependencies and configuring project"
	go mod download && go mod tidy
	@if [ ! -f .env ]; then \
		echo "No .env file found! Please create one based on .env.example."; \
		exit 1; \
	fi
	@echo ".env file found. Proceeding with setup..."
	@echo "Project set up!"

css:
	@echo "Minifying css"
	./tailwindcss -i static/styles/input.css -o static/styles/output.css --minify

css-watch:
	@echo "Watching HTML files for changes"
	./tailwindcss -i static/styles/input.css -o static/styles/output.css --watch

run:
	@echo "Starting PostgreSQL container..."
	@docker compose up -d db
	@echo "Waiting for PostgreSQL to be ready..."
	@until docker exec -it $(shell docker compose ps -q db) pg_isready -U dev_user -d dev_db >/dev/null 2>&1; do \
		echo "Waiting for the database..."; \
		sleep 1; \
	done
	@echo "Database is ready!"
	@echo "Running Go application..."
	go run .

stop-db:
	docker compose down

migrate-down:
	go run data/sql/migrations/migrateDown.go