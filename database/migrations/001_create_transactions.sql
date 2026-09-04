-- Create the transactions table.
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    amount NUMERIC(15, 2) NOT NULL,
    category TEXT NOT NULL,
    type TEXT NOT NULL
        CHECK (type IN ('income', 'expense'))
);