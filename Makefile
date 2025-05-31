# Makefile for prism-user-service
#
# Variables
APP_NAME := prism-user-service
CMD_DIR := cmd/server
MAIN_FILE := $(CMD_DIR)/main.go
BINARY := $(APP_NAME)
GO := go
GOFLAGS := -v
MIGRATE := migrate
MIGRATION_DIR := migrations
DOCKER_COMPOSE := docker-compose
DOCKER_COMPOSE_TEST := $(DOCKER_COMPOSE) -f docker-compose.test.yml
DEV_DB_HOST := db
DEV_DB_PORT := 5432
DEV_DB_USER := prism
DEV_DB_PASSWORD := prism123
DEV_DB_NAME := prism_erp
TEST_DB_HOST := localhost
TEST_DB_PORT := 5433
TEST_DB_USER := test_user
TEST_DB_PASSWORD := test_password
TEST_DB_NAME := test_db
TEST_REDIS_HOST := localhost
TEST_REDIS_PORT := 6380
TEST_ENV := TEST_DATABASE_HOST=$(TEST_DB_HOST) TEST_DATABASE_PORT=$(TEST_DB_PORT) TEST_DATABASE_USER=$(TEST_DB_USER) TEST_DATABASE_PASSWORD=$(TEST_DB_PASSWORD) TEST_DATABASE_NAME=$(TEST_DB_NAME) TEST_REDIS_HOST=$(TEST_REDIS_HOST) TEST_REDIS_PORT=$(TEST_REDIS_PORT)
DEV_DSN := postgres://$(DEV_DB_USER):$(DEV_DB_PASSWORD)@$(DEV_DB_HOST):$(DEV_DB_PORT)/$(DEV_DB_NAME)?sslmode=disable
TEST_DSN := postgres://$(TEST_DB_USER):$(TEST_DB_PASSWORD)@$(TEST_DB_HOST):$(TEST_DB_PORT)/$(TEST_DB_NAME)?sslmode=disable

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	@echo "Building $(APP_NAME)..."
	$(GO) build $(GOFLAGS) -o $(BINARY) $(MAIN_FILE)

# Run the application locally
.PHONY: run
run: build
	@echo "Running $(APP_NAME)..."
	./$(BINARY)

# Build Docker image
.PHONY: docker-build
docker-build:
	@echo "Building Docker image for $(APP_NAME)..."
	$(DOCKER_COMPOSE) build

# Start Docker Compose services
.PHONY: docker-up
docker-up:
	@echo "Starting Docker Compose services..."
	$(DOCKER_COMPOSE) up -d
	@sleep 5
	@echo "Applying database migrations..."
	$(MAKE) migrate-up-dev

# Stop Docker Compose services
.PHONY: docker-down
docker-down:
	@echo "Stopping Docker Compose services..."
	$(DOCKER_COMPOSE) down

# Run tests (unit and integration)
.PHONY: test
test: test-unit test-integration

# Run unit tests
.PHONY: test-unit
test-unit:
	@echo "Running unit tests..."
	$(TEST_ENV) $(GO) test $(GOFLAGS) ./internal/...

# Run integration tests
.PHONY: test-integration
test-integration:
	@echo "Starting test environment..."
	$(DOCKER_COMPOSE_TEST) up -d
	@sleep 10
	@echo "Applying test database migrations..."
	$(MAKE) migrate-up-test
	@echo "Running integration tests..."
	$(TEST_ENV) $(GO) test $(GOFLAGS) -tags=integration ./internal/...
	@echo "Stopping test environment..."
	$(DOCKER_COMPOSE_TEST) down

# Generate test coverage report
.PHONY: coverage
coverage:
	@echo "Generating coverage report..."
	$(TEST_ENV) $(GO) test -coverprofile=coverage.out ./internal/...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run linter
.PHONY: lint
lint:
	@echo "Running linter..."
	golangci-lint run ./...

# Apply database migrations (development)
.PHONY: migrate-up-dev
migrate-up-dev:
	@echo "Applying development database migrations..."
	$(MIGRATE) -path $(MIGRATION_DIR) -database "$(DEV_DSN)" up

# Roll back database migrations (development)
.PHONY: migrate-down-dev
migrate-down-dev:
	@echo "Rolling back development database migrations..."
	$(MIGRATE) -path $(MIGRATION_DIR) -database "$(DEV_DSN)" down

# Apply database migrations (test)
.PHONY: migrate-up-test
migrate-up-test:
	@echo "Applying test database migrations..."
	$(MIGRATE) -path $(MIGRATION_DIR) -database "$(TEST_DSN)" up

# Roll back database migrations (test)
.PHONY: migrate-down-test
migrate-down-test:
	@echo "Rolling back test database migrations..."
	$(MIGRATE) -path $(MIGRATION_DIR) -database "$(TEST_DSN)" down

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning up..."
	rm -f $(BINARY) coverage.out coverage.html

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	$(GO) mod tidy
	$(GO) mod download

# Start test environment
.PHONY: test-env-up
test-env-up:
	@echo "Starting test environment..."
	$(DOCKER_COMPOSE_TEST) up -d

# Stop test environment
.PHONY: test-env-down
test-env-down:
	@echo "Stopping test environment..."
	$(DOCKER_COMPOSE_TEST) down

# Help
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make all               - Build the application"
	@echo "  make build            - Build the application binary"
	@echo "  make run              - Build and run the application locally"
	@echo "  make docker-build     - Build Docker image"
	@echo "  make docker-up        - Start Docker Compose services (applies migrations)"
	@echo "  make docker-down      - Stop Docker Compose services"
	@echo "  make test             - Run unit and integration tests"
	@echo "  make test-unit        - Run unit tests"
	@echo "  make test-integration - Run integration tests with Docker Compose"
	@echo "  make coverage         - Generate test coverage report"
	@echo "  make lint             - Run linter"
	@echo "  make migrate-up-dev   - Apply development database migrations"
	@echo "  make migrate-down-dev - Roll back development database migrations"
	@echo "  make migrate-up-test  - Apply test database migrations"
	@echo "  make migrate-down-test - Roll back test database migrations"
	@echo "  make clean            - Remove build artifacts"
	@echo "  make deps             - Install dependencies"
	@echo "  make test-env-up      - Start test environment (Docker Compose)"
	@echo "  make test-env-down    - Stop test environment (Docker Compose)"
	@echo "  make help             - Show this help message"
