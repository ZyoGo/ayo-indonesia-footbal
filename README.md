# Ayo Indonesia Football Management System

A robust, scalable backend service for managing football clubs, matches, and competition statistics. Built with Go, PostgreSQL, and designed using **Domain-Driven Design (DDD)** and **Clean Architecture** principles.

## Features

*   **Authentication**: JWT-based auth with login/register. Protects write operations.
*   **Club Management**: Register teams and manage player rosters. Protects against duplicate jersey numbers within a team.
*   **Match Management**: Schedule matches between teams, ensuring valid times and no double-booking. Report match results and individual player goals with strict validation (ensuring goal counts match the final score).
*   **Reporting & Analytics**: Automatically aggregates match results into real-time standings (klasemen) based on Points, Goal Difference, and Goals For. Tracks top goalscorers across the competition.

## Documentation

*   [System Architecture & Design Principles](ARCHITECTURE.md)
*   [Entity Relationship Diagram (ERD)](erd.md)

## Tech Stack

*   **Language**: Go 1.26
*   **Framework**: Gin (HTTP Web Framework)
*   **Database**: PostgreSQL 16
*   **Driver/Query Builder**: pgx/v5
*   **IDs**: ULID (Universally Unique Lexicographically Sortable Identifier)
*   **Tooling**: Docker, Docker Compose, Makefile, Air (Hot Reloading)

## Prerequisites

To run this project, you need:
*   [Docker](https://docs.docker.com/get-docker/) & [Docker Compose](https://docs.docker.com/compose/install/)
*   [Make](https://www.gnu.org/software/make/) (Optional, but recommended for ease of use)
*   [Go 1.26+](https://go.dev/dl/) (If running locally without Docker)

## Getting Started

### 1. Running with Docker (Recommended for Development)

We provide a `docker-compose` setup that includes the Postgres database and the Go application with hot-reloading (via `air`).

```bash
# Start the application and database in the background
make up

# Check the logs to ensure everything is running securely
make logs
```

The API will be available at `http://localhost:4000`.

### 2. Database Migrations

The database schema is managed via standard SQL files located in the `migrations/` directory.

```bash
# Run all UP migrations
make migrate-up

# Rollback migrations (DOWN)
make migrate-down
```

### 3. Running Tests

The project includes comprehensive unit tests across the domain and application layers using table-driven tests and `go.uber.org/mock`.

```bash
# Run all tests
make test
```

## API Modules Overview

All endpoints are prefixed with `/api/v1`.

### Authentication (`/auth`)
*   `POST /auth/register`: Register a new user.
*   `POST /auth/login`: Login and receive JWT token.

### Club Context (`/teams`, `/players`)
*   `POST /teams`: Register a new team (protected).
*   `GET /teams`: List all teams.
*   `GET /teams/:id`: Get team by ID.
*   `PUT /teams/:id`: Update team (protected).
*   `DELETE /teams/:id`: Delete team (protected).
*   `POST /players`: Add a player to a team (protected).
*   `GET /players/:id`: Get player by ID.
*   `GET /teams/:id/players`: List all players in a team.
*   `PUT /players/:id`: Update player (protected).
*   `DELETE /players/:id`: Delete player (protected).

### Match Context (`/matches`)
*   `POST /matches`: Schedule a new match (protected).
*   `GET /matches`: List all matches.
*   `GET /matches/:id`: Get match by ID.
*   `GET /matches/:id/report`: Get a detailed report for a specific match.
*   `POST /matches/:id/result`: Report the final result and goal scorers for a match (protected).

### Reporting Context (`/reporting`)
*   `GET /reporting/standings`: Get the current competition standings (klasemen).
*   `GET /reporting/top-scorers`: Get the top goalscorers leaderboard.

### Upload (`/uploads`)
*   `POST /uploads`: Upload a file (protected).

## Project Structure

```text
├── cmd/http/          # Application entrypoint (main.go) and dependency wiring
├── config/            # Application configuration
├── internal/          # Bounded Contexts (club, match, reporting)
│   ├── app/           # Application Services (Use Cases)
│   ├── domain/        # Core Business Logic & Entity definitions
│   └── infra/         # External integrations (Postgres, HTTP Handlers)
├── pkg/               # Shared utilities (errors, http responses, logger)
└── migrations/        # SQL migration files
```

## Stopping the Application

```bash
make down
```