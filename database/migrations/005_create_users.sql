-- Migration 005
-- Create the users table.
--
-- Each user will eventually have their own
-- transactions, categories, and account.

BEGIN;


-- Create the users table.
CREATE TABLE users (

    -- Unique identifier for each user.
    id SERIAL PRIMARY KEY,

    -- User's login email.
    --
    -- UNIQUE prevents two accounts from
    -- registering with the same email.
    email TEXT NOT NULL UNIQUE,

    -- Securely hashed password.
    --
    -- We will NEVER store the user's
    -- plain-text password here.
    password_hash TEXT NOT NULL,

    -- When the user account was created.
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()

);


-- Create an index on email.
--
-- PostgreSQL already creates an index for the
-- UNIQUE constraint, so we do not need another
-- separate email index.


COMMIT;
