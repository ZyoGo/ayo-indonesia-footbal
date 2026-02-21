# Entity Relationship Diagram (ERD)

This document contains the Entity Relationship Diagram for the Football Management System, generated using Mermaid.

```mermaid
erDiagram
    users {
        varchar(26) id PK "ULID"
        varchar(255) username "UNIQUE"
        text password_hash
        timestamptz created_at
    }

    teams {
        varchar(26) id PK "ULID"
        varchar(255) name
        text logo_url
        integer year_founded
        text address
        varchar(100) city
        timestamptz created_at
        timestamptz updated_at
        timestamptz deleted_at "Soft Delete"
    }

    players {
        varchar(26) id PK "ULID"
        varchar(26) team_id FK
        varchar(255) name
        decimal(5_2) height
        decimal(5_2) weight
        varchar(20) position
        integer jersey_number
        timestamptz created_at
        timestamptz updated_at
        timestamptz deleted_at "Soft Delete"
    }

    matches {
        varchar(26) id PK "ULID"
        varchar(26) home_team_id FK
        varchar(26) away_team_id FK
        date match_date
        varchar(5) match_time "HH:MM"
        varchar(255) stadium
        timestamptz created_at
        timestamptz updated_at
        timestamptz deleted_at "Soft Delete"
    }

    match_results {
        varchar(26) id PK "ULID"
        varchar(26) match_id FK "UNIQUE"
        integer home_score
        integer away_score
        timestamptz deleted_at "Soft Delete"
    }

    goals {
        varchar(26) id PK "ULID"
        varchar(26) result_id FK
        varchar(26) player_id FK
        varchar(26) team_id FK
        integer goal_minute
        timestamptz deleted_at "Soft Delete"
    }

    %% Relationships
    users ||--o{ "": ""
    teams ||--o{ players : "has"
    teams ||--o{ matches : "plays as home"
    teams ||--o{ matches : "plays as away"
    teams ||--o{ goals : "scores"
    
    matches ||--o| match_results : "has result"
    
    match_results ||--o{ goals : "includes"
    
    players ||--o{ goals : "scores"
```

## Description of Entities

*   **`users`**: Stores user credentials for JWT-based authentication.
*   **`teams`**: Represents a football club.
*   **`players`**: Represents a football player who belongs to a `team`. A team cannot have two players with the same jersey number (enforced by a composite unique constraint).
*   **`matches`**: Represents a scheduled game between a home team and an away team.
*   **`match_results`**: Stores the final score of a `match`. It has a strict 1-to-1 relationship with `matches` (via a unique constraint on `match_id`).
*   **`goals`**: Records an individual goal scored during a match result. It points to the `match_result` it belongs to, the `player` who scored it, and the `team` the player scored for.
