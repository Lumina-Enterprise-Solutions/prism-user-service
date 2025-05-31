#!/bin/bash

set -e

echo "Starting test environment..."

# Start test dependencies
docker-compose -f docker-compose.test.yml up -d

# Wait for services to be ready
echo "Waiting for services to be ready..."
sleep 10

# Set test environment variables
export TEST_DATABASE_HOST=localhost
export TEST_DATABASE_PORT=5433
export TEST_DATABASE_USER=test_user
export TEST_DATABASE_PASSWORD=test_password
export TEST_DATABASE_NAME=test_db
export TEST_REDIS_HOST=localhost
export TEST_REDIS_PORT=6380

# Run tests
echo "Running unit tests..."
go test -v ./internal/...

echo "Running integration tests..."
go test -v -tags=integration ./internal/...

# Generate coverage report
echo "Generating coverage report..."
go test -coverprofile=coverage.out ./internal/...
go tool cover -html=coverage.out -o coverage.html

echo "Stopping test environment..."
docker-compose -f docker-compose.test.yml down

echo "Tests completed successfully!"
