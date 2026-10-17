# Go API with Fiber

A REST API built with **Golang and Fiber**, featuring PostgreSQL database integration, Redis connectivity, JWT-based authentication support, and request rate limiting.

This project explores backend development with Go, including API routing, database connections, authentication, middleware, and environment-based configuration.

## Tech Stack

- **Language:** Go
- **Web Framework:** [Fiber v2](https://github.com/gofiber/fiber)
- **Database:** PostgreSQL
- **ORM:** GORM
- **Cache / Data Store:** Redis
- **Authentication:** JSON Web Tokens (JWT)
- **Configuration:** Environment variables with `godotenv`
- **Middleware:** Fiber rate limiter

## Features

- REST API routing with Fiber
- PostgreSQL integration using GORM
- Redis connection support
- JWT access, refresh, and anonymous token configuration
- Request rate limiting
- Environment-based configuration for credentials and connection settings
- Modular code organization for API handlers, database connections, middleware, and utilities

## Project Structure

```text
go-api-with-fiber/
├── api/           # API routes and handlers
├── database/      # Database and Redis connections
├── middleware/    # HTTP middleware
├── utils/         # Shared utility functions
├── .env.sample    # Environment variable template
├── .gitignore
├── go.mod
├── go.sum
└── main.go        # Application entry point
```

## Prerequisites

Make sure you have the following installed:

- [Go](https://go.dev/dl/) 1.22.4 or a compatible version
- [PostgreSQL](https://www.postgresql.org/)
- [Redis](https://redis.io/)

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/Permanaa/go-api-with-fiber.git
cd go-api-with-fiber
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Configure environment variables

Copy the sample environment file:

```bash
cp .env.sample .env
```

On Windows PowerShell, use:

```powershell
Copy-Item .env.sample .env
```

Update `.env` with your local configuration:

```env
# JWT
JWT_ACCESS_TOKEN_SECRET_KEY=your-access-token-secret
JWT_REFRESH_TOKEN_SECRET_KEY=your-refresh-token-secret
JWT_ANONYMOUS_TOKEN_SECRET_KEY=your-anonymous-token-secret

# PostgreSQL
DB_USER=postgres
DB_PASSWORD=your-database-password
DB_HOST=localhost
DB_PORT=5432
DB_NAME=postgres

# Basic Authentication
BASIC_AUTH_USERNAME=your-username
BASIC_AUTH_PASSWORD=your-password

# Redis
REDIS_ADDRESS=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
```

Use strong, unique secrets for JWT keys and never commit your `.env` file or production credentials.

### 4. Start PostgreSQL and Redis

Ensure both services are running and that the database credentials match your environment configuration.

### 5. Run the application

```bash
go run .
```

The application is configured to listen on:

```text
http://127.0.0.1:1804
```

## API Usage

The available endpoints are defined in the `api/` directory.

To inspect the registered routes, review the route configuration and handlers in that directory.

You can test the API using tools such as:

- [Postman](https://www.postman.com/)
- [Insomnia](https://insomnia.rest/)
- `curl`

Example request:

```bash
curl -i http://127.0.0.1:1804/
```

Replace the path with an endpoint registered by the application.

## Running Tests

Run all Go tests:

```bash
go test ./...
```

Check the project for compilation issues:

```bash
go build ./...
```

## Learning Objectives

This project provides practical experience with:

- Building HTTP APIs in Go
- Organizing backend code into separate packages
- Connecting an application to PostgreSQL and Redis
- Managing configuration through environment variables
- Working with JWT authentication
- Applying middleware to HTTP requests
- Developing backend services using Go and Fiber
