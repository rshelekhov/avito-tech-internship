# Run linters
lint:
	@echo "Running linters..."
	golangci-lint run --fix
	@echo "Linters completed."