-- Rollback: Drop Club Management tables

DROP INDEX IF EXISTS idx_teams_active;
DROP INDEX IF EXISTS idx_players_team_id;
DROP INDEX IF EXISTS idx_unique_jersey_per_team;
DROP TABLE IF EXISTS players;
DROP TABLE IF EXISTS teams;
