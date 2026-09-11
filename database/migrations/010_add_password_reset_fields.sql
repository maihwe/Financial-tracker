-- Migration 010
-- Add password-reset fields to the users table.

BEGIN;

ALTER TABLE users
ADD COLUMN password_reset_token TEXT UNIQUE,
ADD COLUMN password_reset_expires_at TIMESTAMPTZ;

COMMIT;