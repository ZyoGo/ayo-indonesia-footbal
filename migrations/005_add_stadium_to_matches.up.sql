-- Migration: Add stadium to matches
ALTER TABLE matches ADD COLUMN IF NOT EXISTS stadium VARCHAR(255);
