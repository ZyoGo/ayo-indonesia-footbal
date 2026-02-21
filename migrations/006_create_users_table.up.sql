-- Migration: Create users table for authentication
-- Description: Stores admin/user credentials for JWT-based authentication

CREATE TABLE IF NOT EXISTS users (
    id            VARCHAR(26) PRIMARY KEY,
    username      VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
