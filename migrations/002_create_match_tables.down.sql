-- Rollback: Drop Match & Competition tables

DROP INDEX IF EXISTS idx_goals_player_id;
DROP INDEX IF EXISTS idx_goals_result_id;
DROP INDEX IF EXISTS idx_matches_date;
DROP TABLE IF EXISTS goals;
DROP TABLE IF EXISTS match_results;
DROP TABLE IF EXISTS matches;
