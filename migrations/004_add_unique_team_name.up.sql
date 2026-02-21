-- Migration: Add unique constraint to team name
-- Description: Ensures that active teams cannot have the exact same name

CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_team_name
    ON teams (name)
    WHERE deleted_at IS NULL;
