# Entity Relationship Diagram (ERD)

This document contains the Entity Relationship Diagram for the Football Management System, generated using Mermaid.

```mermaid
erDiagram
    teams {
        char(26) id PK "ULID"
        varchar(100) name "UNIQUE"
        varchar(100) city
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at "Soft Delete"
    }

    players {
        char(26) id PK "ULID"
        char(26) team_id FK
        varchar(100) name
        numeric(5_2) height
        numeric(5_2) weight
        varchar(20) position
        int jersey_number
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at "Soft Delete"
    }

    matches {
        char(26) id PK "ULID"
        char(26) home_team_id FK
        char(26) away_team_id FK
        date match_date
        varchar(5) match_time "HH:MM"
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at "Soft Delete"
    }

    match_results {
        char(26) id PK "ULID"
        char(26) match_id FK "UNIQUE"
        int home_score
        int away_score
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at "Soft Delete"
    }

    goals {
        char(26) id PK "ULID"
        char(26) result_id FK
        char(26) player_id FK
        char(26) team_id FK
        int goal_minute
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at "Soft Delete"
    }

    %% Relationships
    teams ||--o{ players : "has"
    teams ||--o{ matches : "plays as home"
    teams ||--o{ matches : "plays as away"
    teams ||--o{ goals : "scores"
    
    matches ||--o| match_results : "has result"
    
    match_results ||--o{ goals : "includes"
    
    players ||--o{ goals : "scores"
```

## Description of Entities

*   **`teams`**: Represents a football club.
*   **`players`**: Represents a football player who belongs to a `team`. A team cannot have two players with the same jersey number (enforced by a composite unique constraint).
*   **`matches`**: Represents a scheduled game between a home team and an away team.
*   **`match_results`**: Stores the final score of a `match`. It has a strict 1-to-1 relationship with `matches` (via a unique constraint on `match_id`).
*   **`goals`**: Records an individual goal scored during a match result. It points to the `match_result` it belongs to, the `player` who scored it, and the `team` the player scored for.
