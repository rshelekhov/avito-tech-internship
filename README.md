# Merch Store Service

A backend service for an internal merchandise store where employees can purchase items using coins and transfer coins between each other. This implementation is based on the [Winter 2025 Backend Internship test assignment](https://github.com/avito-tech/tech-internship/blob/main/Tech%20Internships/Backend/Backend-trainee-assignment-winter-2025/Backend-trainee-assignment-winter-2025.md) by Avito.

## Features

- User authentication with JWT
- Coin transfer between employees
- Merchandise purchase system
- Transaction history tracking
- Automatic new user registration with 1000 coins initial balance
- Comprehensive test coverage with unit and E2E tests

## Built With

- PostgreSQL as the main database
- [golang-migrate](https://github.com/golang-migrate/migrate) for database migrations
- [sqlc](https://github.com/sqlc-dev/sqlc) for generating type-safe code from SQL
- [viper](https://github.com/spf13/viper) for configuration management
- [log/slog](https://pkg.go.dev/log/slog) for centralized logging
- [golangci-lint](https://github.com/golangci/golangci-lint) for code quality
- [mockery](https://github.com/vektra/mockery) for mock generation

## Project Structure

The project follows a clean architecture approach with the following key components:

- Domain entities and services
- Use cases implementing business logic
- HTTP handlers and middleware
- Database migrations and storage layer
- Configuration management
- Comprehensive testing suite

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Make (optional, but recommended)

### Running the Application

1. Clone the repository:
```bash
git clone https://github.com/rshelekhov/merch-store.git
cd merch-store
```

2. Start the application:
```bash
docker compose up postgres app
```

The service will be available at `http://localhost:8080`

### Running Tests

To run all tests (unit and E2E):
```bash
docker compose run --rm test
```


### Configuration

The application uses environment variables for configuration. You have two options:

- Provide your own configuration:
    - Create a `.env` file in the `./config` directory
    - Configure the environment variables according to your needs

- Use default configuration:
  - If no `.env` file is present, the container will automatically copy `.env.example` to `.env`
  - The application will start with these default settings

The configuration file should be mounted to `/src/config/.env` in the container.

### Local Development
For local development without containers:
- Use local.env configuration file
- Ensure PostgreSQL is running locally
- Create a database for the application before starting

## License

This project is licensed under the MIT License - see the LICENSE file for details.
