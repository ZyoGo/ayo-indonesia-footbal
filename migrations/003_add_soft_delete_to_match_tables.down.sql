-- Rollback: Remove soft delete support from Match & Competition tables

DROP INDEX IF EXISTS idx_matches_active;
DROP INDEX IF EXISTS idx_match_results_active;
DROP INDEX IF EXISTS idx_goals_active;

ALTER TABLE matches DROP COLUMN deleted_at;
ALTER TABLE match_results DROP COLUMN deleted_at;
ALTER TABLE goals DROP COLUMN deleted_at;
