-- Migration: Remove stadium from matches
ALTER TABLE matches DROP COLUMN IF EXISTS stadium;
