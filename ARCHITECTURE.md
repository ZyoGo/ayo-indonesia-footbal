# System Architecture

This document describes the architectural principles, patterns, and structure used in the Football Management System.

## Core Architectural Principles

The application is built using two primary architectural patterns:
1.  **Domain-Driven Design (DDD)**: Organizing the system around business domains and bounded contexts.
2.  **Clean Architecture (Hexagonal / Ports and Adapters)**: Separating technical concerns from business logic via strict dependency inversion.

### Bounded Contexts (DDD)

We have divided the system into four distinct bounded contexts to ensure high cohesion and loose coupling:

1.  **Auth Context**: Manages user authentication and JWT token generation.
2.  **Club Context**: Manages the core entities of football teams and their players. Responsible for team registration and squad management.
3.  **Match Context**: Manages match scheduling, storing results, and tracking individual match events (e.g., goals scored, minutes played). Focuses purely on the transactional aspect of playing a game.
4.  **Reporting Context**: A read-heavy context responsible for aggregating data from matches to generate standings (klasemen) and player leaderboards (top scorers). It observes match results but does not manage them directly.

---

## Clean Architecture Layers

Each bounded context (inside the `internal/` directory) strictly follows Clean Architecture layers:

### 1. Domain Layer (`domain/`)
-   **Definition**: The heart of the software. Contains enterprise-wide business rules, entities, and repository interfaces (Ports).
-   **Dependencies**: **None**. It must not import anything from `app`, `infra`, or standard libraries tied to frameworks (like `gin` or `pgx`).
-   **Role**: Defines *what* the system does (e.g., a `Team` cannot have duplicate names, a `MatchResult` must have goals matching the score).

### 2. Application Layer (`app/`)
-   **Definition**: Contains use cases and application-specific business rules.
-   **Dependencies**: Depends *only* on the `domain` layer.
-   **Role**: Orchestrates the flow of data. It uses the repository interfaces (Ports) defined in the domain to fetch/save data, orchestrates validation, and returns results.

### 3. Infrastructure Layer (`infra/`)
-   **Definition**: The outermost layer containing framework integrations, database drivers, and HTTP handlers.
-   **Dependencies**: Depends on both `domain` and `app` layers.
-   **Role**: Implements the technical details (Adapters).
    -   `postgres/`: Implements repository interfaces using PostgreSQL and `pgx`.
    -   `handler/`: Implements HTTP controllers using `gin`, mapping HTTP requests to Service calls.

---

## Directory Structure

```text
├── cmd/
│   └── http/                 # Application entrypoint (main.go) and module wiring
├── config/                   # Configuration binding and defaults (Viper)
├── internal/                 # Private application code (Bounded Contexts)
│   ├── auth/                 # Auth Context
│   │   ├── app/              # Application logic (Services)
│   │   ├── domain/           # Core rules and interfaces
│   │   ├── infra/           # Postgres implementations, HTTP Handlers
│   │   └── mock/             # Generated mocks for testing
│   ├── club/                 # Club Context
│   │   ├── app/              # Application logic (Services)
│   │   ├── domain/           # Core rules and interfaces
│   │   ├── infra/            # Postgres implementations, HTTP Handlers
│   │   └── mock/             # Generated mocks for testing
│   ├── match/                # Match Context
│   │   └── ...
│   └── reporting/            # Reporting Context
│       └── ...
├── pkg/                      # Shared, generic utilities (logging, http responses, custom errors, upload, jwt)
├── migrations/               # Raw SQL files for database migrations
└── Dockerfile                # Multi-stage build definitions
```

---

## Data Flow (Request Lifecycle)

A typical HTTP request follows this exact path (Dependency Inversion in action):

1.  **Incoming Request**: Client makes an HTTP request to an endpoint (e.g., `POST /api/v1/clubs`).
2.  **Infrastructure (Handler)**: The Gin Router directs the request to a specific `Handler`. The handler parses the JSON `Request` body, checks basic syntactical validity, and calls a method on the `Service`.
3.  **Application (Service)**: The `Service` receives the request. It holds business orchestration logic. Before saving, it may call business rules defined on the `Domain` entity. To save data, it calls a method on its injected `Repository Interface`.
4.  **Infrastructure (Database)**: The injected `PostgresRepository` receives the domain entity, maps it to SQL statements, and executes the queries via `pgx`.
5.  **Response**: The flow reverses. The database returns rows, the repository returns domain entities, the service returns data to the handler, and the handler maps it back to a standard JSON `Response`.

## Error Handling Pattern

We use a custom, centralized error package `pkg/derrors` that allows propagating HTTP status codes (400, 404, 500) natively from the Application layer, ensuring Handlers don't leak domain implementation details while still returning precise HTTP responses.
