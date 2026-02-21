-- Migration: Add soft delete support to Match & Competition tables
-- Description: Adds deleted_at columns and indexes for soft delete

ALTER TABLE matches ADD COLUMN deleted_at TIMESTAMPTZ;
ALTER TABLE match_results ADD COLUMN deleted_at TIMESTAMPTZ;
ALTER TABLE goals ADD COLUMN deleted_at TIMESTAMPTZ;

-- Indexes for faster active record filtering
CREATE INDEX IF NOT EXISTS idx_matches_active ON matches (id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_match_results_active ON match_results (id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_goals_active ON goals (id) WHERE deleted_at IS NULL;
