-- Migration: Create tables for Club Management Context
-- Description: Creates teams and players tables with soft delete support

CREATE TABLE IF NOT EXISTS teams (
    id          VARCHAR(26) PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    logo_url    TEXT DEFAULT '',
    year_founded INTEGER NOT NULL,
    address     TEXT DEFAULT '',
    city        VARCHAR(100) NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS players (
    id              VARCHAR(26) PRIMARY KEY,
    team_id         VARCHAR(26) NOT NULL REFERENCES teams(id),
    name            VARCHAR(255) NOT NULL,
    height          DECIMAL(5,2) NOT NULL DEFAULT 0,
    weight          DECIMAL(5,2) NOT NULL DEFAULT 0,
    position        VARCHAR(20) NOT NULL,
    jersey_number   INTEGER NOT NULL CHECK (jersey_number BETWEEN 1 AND 99),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

-- Partial unique index: ensures jersey number uniqueness per team among active (non-deleted) players
CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_jersey_per_team
    ON players (team_id, jersey_number)
    WHERE deleted_at IS NULL;

-- Index for faster team lookups on players
CREATE INDEX IF NOT EXISTS idx_players_team_id ON players (team_id) WHERE deleted_at IS NULL;

-- Index for soft delete filtering on teams
CREATE INDEX IF NOT EXISTS idx_teams_active ON teams (id) WHERE deleted_at IS NULL;
