-- Migration 006
-- Add user ownership to transactions.
--
-- We temporarily allow NULL because transactions
-- already exist in the database and do not yet have
-- an owner.

BEGIN;


-- Add the user_id column.
--
-- NULL is temporarily allowed while we migrate
-- existing transactions to an owner.
ALTER TABLE transactions

ADD COLUMN user_id INTEGER;


-- Connect transactions to the users table.
--
-- This ensures that when user_id contains a value,
-- that value must exist in users.id.
ALTER TABLE transactions

ADD CONSTRAINT fk_transactions_user

FOREIGN KEY (user_id)

REFERENCES users(id);


COMMIT;
