-- Migration: Create tables for Match & Competition Context
-- Description: Creates matches, match_results, and goals tables

CREATE TABLE IF NOT EXISTS matches (
    id              VARCHAR(26) PRIMARY KEY,
    home_team_id    VARCHAR(26) NOT NULL REFERENCES teams(id),
    away_team_id    VARCHAR(26) NOT NULL REFERENCES teams(id),
    match_date      DATE NOT NULL,
    match_time      VARCHAR(5) NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_different_teams CHECK (home_team_id != away_team_id)
);

CREATE TABLE IF NOT EXISTS match_results (
    id              VARCHAR(26) PRIMARY KEY,
    match_id        VARCHAR(26) NOT NULL UNIQUE REFERENCES matches(id),
    home_score      INTEGER NOT NULL DEFAULT 0 CHECK (home_score >= 0),
    away_score      INTEGER NOT NULL DEFAULT 0 CHECK (away_score >= 0)
);

CREATE TABLE IF NOT EXISTS goals (
    id              VARCHAR(26) PRIMARY KEY,
    result_id       VARCHAR(26) NOT NULL REFERENCES match_results(id),
    player_id       VARCHAR(26) NOT NULL REFERENCES players(id),
    team_id         VARCHAR(26) NOT NULL REFERENCES teams(id),
    goal_minute     INTEGER NOT NULL CHECK (goal_minute > 0)
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_matches_date ON matches (match_date DESC);
CREATE INDEX IF NOT EXISTS idx_goals_result_id ON goals (result_id);
CREATE INDEX IF NOT EXISTS idx_goals_player_id ON goals (player_id);
