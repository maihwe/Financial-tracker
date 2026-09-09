-- Migration 004
-- Align existing transaction_at values with created_at.
--
-- These transactions existed before transaction_at was added.
-- created_at is therefore the best-known timestamp for them.
--
-- This migration runs only once.
-- Future transactions will continue using transaction_at
-- DEFAULT NOW() from Migration 003.

BEGIN;

UPDATE transactions
SET transaction_at = created_at;

DO $$
BEGIN

    IF EXISTS (
        SELECT 1
        FROM transactions
        WHERE transaction_at <> created_at
    ) THEN

        RAISE EXCEPTION
            'Migration stopped: some transaction timestamps were not aligned.';

    END IF;

END $$;

COMMIT;
