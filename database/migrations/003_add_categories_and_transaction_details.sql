-- Migration 003
-- Add categories and additional transaction details.


-- Start a transaction so that all changes succeed together.
BEGIN;


-- Create the categories table.
--
-- user_id = NULL means the category is global.
-- user_id = a specific user ID means the category
-- belongs to that user.

CREATE TABLE categories (

    id SERIAL PRIMARY KEY,

    name TEXT NOT NULL,

    user_id INTEGER,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()

);


-- Add an optional description to transactions.
--
-- The description allows users to provide
-- additional information about a transaction.

ALTER TABLE transactions

ADD COLUMN description TEXT;


-- Add the actual date and time when the transaction happened.
--
-- transaction_at = when the financial activity happened.
-- created_at = when the transaction record was created.
-- updated_at = when the transaction record was last modified.

ALTER TABLE transactions

ADD COLUMN transaction_at TIMESTAMPTZ NOT NULL DEFAULT NOW();


-- Insert the initial global categories.
--
-- Global categories have user_id = NULL.
-- This means every user will be able to use them.

INSERT INTO categories (name, user_id)
VALUES
    ('Food', NULL),
    ('Transportation', NULL),
    ('Salary', NULL),
    ('Business', NULL),
    ('Education', NULL),
    ('Health', NULL),
    ('Bills & Utilities', NULL),
    ('Housing', NULL),
    ('Shopping', NULL),
    ('Entertainment', NULL),
    ('Savings', NULL),
    ('Debt', NULL),
    ('Investment', NULL),
    ('Other', NULL);


-- Add existing categories that were not in the
-- initial global category list.
--
-- These preserve existing transaction records.

INSERT INTO categories (name, user_id)
VALUES
    ('Trade', NULL),
    ('Hygiene', NULL);


-- Add category_id to transactions.
--
-- This temporarily allows NULL values while
-- existing transactions are connected to categories.

ALTER TABLE transactions

ADD COLUMN category_id INTEGER;


-- Connect existing transactions to their new category IDs.
--
-- The old transactions.category value is matched
-- with categories.name.

UPDATE transactions AS t

SET category_id = c.id

FROM categories AS c

WHERE LOWER(t.category) = LOWER(c.name);


-- Verify that every existing transaction
-- has been connected to a category.

DO $$
BEGIN

    IF EXISTS (
        SELECT 1
        FROM transactions
        WHERE category_id IS NULL
    ) THEN

        RAISE EXCEPTION
            'Migration stopped: some transactions have no category match.';

    END IF;

END $$;


-- Every transaction must have a category.

ALTER TABLE transactions

ALTER COLUMN category_id SET NOT NULL;


-- Connect transactions to the categories table.
--
-- This ensures that every category_id in transactions
-- must refer to an existing category.

ALTER TABLE transactions

ADD CONSTRAINT fk_transactions_category

FOREIGN KEY (category_id)

REFERENCES categories(id);


-- Remove the old category text column.
--
-- Transactions now use category_id to reference
-- the categories table.

ALTER TABLE transactions

DROP COLUMN category;


-- Finish the migration.
COMMIT;