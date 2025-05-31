# Prism User Service
[![Go Report Card](https://goreportcard.com/badge/github.com/Lumina-Enterprise-Solutions/prism-user-service)](https://goreportcard.com/report/github.com/Lumina-Enterprise-Solutions/prism-user-service)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![CI/CD Pipeline](https://github.com/Lumina-Enterprise-Solutions/prism-user-service/actions/workflows/ci.yml/badge.svg)](https://github.com/Lumina-Enterprise-Solutions/prism-user-service/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/Go-1.20+-00ADD8.svg)](https://golang.org/dl/)

## Overview

The **Prism User Service** is a microservice built with Go, designed to manage user-related operations such as user creation, authentication, role management, and profile updates within the Prism ERP system. It leverages PostgreSQL for persistent storage, Redis for caching, and integrates with the `prism-common-libs` package for shared functionality like database connections, middleware, and logging.

This service provides a RESTful API for managing users, roles, and permissions, with support for multi-tenancy, JWT-based authentication, and comprehensive logging.

## Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [Directory Structure](#directory-structure)
- [Prerequisites](#prerequisites)
- [Setup](#setup)
- [Running the Service](#running-the-service)
- [Testing](#testing)
- [API Endpoints](#api-endpoints)
- [Environment Variables](#environment-variables)
- [Contributing](#contributing)
- [License](#license)

## Features

- **User Management**: Create, update, delete, and list users with pagination and filtering.
- **Role-Based Access Control (RBAC)**: Assign roles to users with customizable permissions stored in JSONB.
- **Multi-Tenancy**: Supports tenant-specific data isolation using PostgreSQL schemas.
- **JWT Authentication**: Secure endpoints with token-based authentication.
- **Health Checks**: Provides `/health` and `/ready` endpoints for service monitoring.
- **Logging**: Structured logging with customizable levels and formats (JSON or human-readable).
- **Database Migrations**: Manages database schema changes using SQL migration files.
- **Testing**: Includes unit and integration tests with coverage reporting.

## Architecture

The service follows a clean architecture pattern, separating concerns into layers:

- **Handlers**: Handle HTTP requests and responses using the Gin framework.
- **Services**: Contain business logic for user operations.
- **Repository**: Manage database interactions using GORM.
- **Models**: Define data structures for users, roles, and requests/responses.
- **Config**: Manage configuration loading from environment variables.
- **Middleware**: Handle cross-cutting concerns like authentication, CORS, and request IDs.

The service integrates with PostgreSQL for data storage and Redis for caching, using the `prism-common-libs` package for shared utilities.

## Directory Structure

```
prism-user-service/
├── cmd/
│   └── server/
│       └── main.go                # Entry point for the service
├── internal/
│   ├── config/                    # Configuration loading
│   │   └── config.go
│   ├── handlers/                  # HTTP handlers
│   │   ├── health.go
│   │   └── user.go
│   ├── models/                    # Data models
│   │   └── user.go
│   ├── repository/                # Database operations
│   │   ├── mock_user_repository.go
│   │   └── user.go
│   └── services/                  # Business logic
│       └── user.go
├── migrations/                    # Database migration scripts
│   ├── 001_create_users_table.up.sql
│   ├── 001_create_users_table.down.sql
│   ├── 002_create_roles_table.up.sql
│   ├── 002_create_roles_table.down.sql
│   ├── 003_create_user_roles_table.up.sql
│   └── 003_create_user_roles_table.down.sql
├── scripts/
│   └── test.sh                    # Script to run tests
├── docker-compose.yml             # Docker Compose configuration
├── .mdignore                      # Files to ignore in documentation
└── prism-user-service             # Binary (not included in source)
```

## Prerequisites

- **Go**: Version 1.16 or higher
- **Docker** and **Docker Compose**: For running PostgreSQL and Redis
- **PostgreSQL**: Version 13 or higher
- **Redis**: Version 6 or higher
- **Make**: For running build and test scripts (optional)

## Setup

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/Lumina-Enterprise-Solutions/prism-user-service.git
   cd prism-user-service
   ```

2. **Install Dependencies**:
   ```bash
   go mod download
   ```

3. **Set Up Environment Variables**:
   Create a `.env` file or set environment variables as described in [Environment Variables](#environment-variables).

4. **Start Dependencies**:
   Start PostgreSQL and Redis using Docker Compose:
   ```bash
   docker-compose up -d
   ```

5. **Apply Database Migrations**:
   Ensure your database migration tool (e.g., `migrate`) is installed, then run:
   ```bash
   migrate -path migrations -database "postgres://prism:prism123@localhost:5432/prism_erp?sslmode=disable" up
   ```

## Running the Service

1. **Build the Service**:
   ```bash
   go build -o prism-user-service ./cmd/server
   ```

2. **Run the Service**:
   ```bash
   ./prism-user-service
   ```

   Alternatively, use Docker Compose:
   ```bash
   docker-compose up app
   ```

   The service will be available at `http://localhost:8080`.

## Testing

Run unit and integration tests using the provided script:
```bash
./scripts/test.sh
```

This script:
- Starts a test environment with Docker Compose
- Sets test environment variables
- Runs unit and integration tests
- Generates a coverage report (`coverage.html`)
- Cleans up the test environment

## API Endpoints

All endpoints are prefixed with `/api/v1`. Protected endpoints require a valid JWT in the `Authorization` header and a tenant ID in the `X-Tenant-ID` header.

| Method | Endpoint                | Description                      | Authentication |
|--------|-------------------------|----------------------------------|----------------|
| GET    | `/health`              | Check service health              | None           |
| GET    | `/ready`               | Check service readiness          | None           |
| POST   | `/users`               | Create a new user                | JWT            |
| GET    | `/users`               | List users with pagination       | JWT            |
| GET    | `/users/:id`           | Get user by ID                   | JWT            |
| PUT    | `/users/:id`           | Update user                      | JWT            |
| DELETE | `/users/:id`           | Delete user                      | JWT            |
| GET    | `/users/profile`       | Get authenticated user's profile  | JWT            |
| PUT    | `/users/profile`       | Update authenticated user's profile | JWT          |

### Example Request
**Create User**:
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -H "X-Tenant-ID: default" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "password": "securepassword123",
    "status": "active",
    "role_ids": ["<ROLE_ID>"]
  }'
```

## Environment Variables

The service uses the following environment variables, with defaults provided:

| Variable                | Description                              | Default Value          |
|-------------------------|------------------------------------------|------------------------|
| `ENVIRONMENT`           | Environment (development/production)     | `development`          |
| `SERVICE_NAME`          | Service name                             | `prism-user-service`   |
| `SERVICE_VERSION`       | Service version                          | `v1.0.0`              |
| `SERVICE_ENVIRONMENT`   | Service environment                      | `development`          |
| `DB_HOST`               | Database host                            | `db`                   |
| `DB_PORT`               | Database port                            | `5432`                |
| `DB_NAME`               | Database name                            | `prism_erp`           |
| `DB_USER`               | Database username                        | `prism`               |
| `DB_PASSWORD`           | Database password                        | `prism123`            |
| `DB_SSL_MODE`           | Database SSL mode                        | `disable`             |
| `REDIS_HOST`            | Redis host                               | `redis`               |
| `REDIS_PORT`            | Redis port                               | `6379`                |
| `REDIS_PASSWORD`        | Redis password                           | `redis123`            |
| `REDIS_DB`              | Redis database number                    | `0`                   |
| `JWT_SECRET`            | JWT secret key                           | `your-secret-key`     |
| `SERVER_HOST`           | Server host                              | `0.0.0.0`             |
| `SERVER_PORT`           | Server port                              | `8080`                |
| `SERVER_READ_TIMEOUT`   | Server read timeout (seconds)            | `10`                  |
| `SERVER_WRITE_TIMEOUT`  | Server write timeout (seconds)           | `10`                  |
| `LOG_LEVEL`             | Log level (debug/info/warn/error/fatal) | `info`                |
| `LOG_FORMAT`            | Log format (json/text)                   | `json`                |

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository.
2. Create a feature branch (`git checkout -b feature/your-feature`).
3. Commit your changes (`git commit -m "Add your feature"`).
4. Push to the branch (`git push origin feature/your-feature`).
5. Open a pull request.

Ensure your code follows the project's coding standards and includes tests.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
