CONFIG_PATH ?= ./config/.env
SERVER_PORT ?= 8080

# Don't forget to set POSTGRESQL_URL with your credentials
POSTGRESQL_URL ?= postgres://root:password@localhost:5432/merch_store_dev?sslmode=disable

# Run migrations only if not already applied
migrate:
	@echo "Checking if postgresql-client is installed..."
	@if ! which psql > /dev/null 2>&1; then \
		echo "postgresql-client not found. Installing..."; \
		if [ "$$(uname)" = "Darwin" ]; then \
			echo "Detected macOS. Installing via Homebrew..."; \
			brew install postgresql; \
		elif [ "$$(uname)" = "Linux" ]; then \
			echo "Detected Linux. Installing via apt-get..."; \
			sudo apt-get update && sudo apt-get install -y postgresql-client; \
		else \
			echo "Unsupported OS. Please install postgresql-client manually."; \
			exit 1; \
		fi \
	else \
		echo "postgresql-client is already installed."; \
	fi

	@echo "Checking if migrations are needed..."
		@if psql $(POSTGRESQL_URL) -c "SELECT 1 FROM pg_tables WHERE tablename = 'apps';" | grep -q 1; then \
			echo "Migrations are not needed."; \
		else \
			echo "Running migrations..."; \
			migrate -database $(POSTGRESQL_URL) -path migrations up; \
			echo "Migrations completed."; \
		fi

# Rollback migrations
migrate-down:
	@echo "Rolling back migrations..."
	@migrate -database $(POSTGRESQL_URL) -path migrations down
	@echo "Migrations rolled back."

# Run linters
lint:
	@echo "Running linters..."
	golangci-lint run --fix
	@echo "Linters completed."