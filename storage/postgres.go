package storage

import (
	"context"

	"financial-tracker/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

// GetAllTransactionsFromDB retrieves every transaction
// from the PostgreSQL database.
func GetAllTransactionsFromDB(
	pool *pgxpool.Pool) ([]models.Transaction, error) {

	// Send a SELECT query to PostgreSQL.
	rows, err := pool.Query(
		context.Background(),
		`
		SELECT id, title, amount, category, type, created_at, updated_at
		FROM transactions
		ORDER BY id
		`,
	)

	// Stop if PostgreSQL returned an error.
	if err != nil {
		return nil, err
	}

	// Make sure the database rows are closed
	// when this function finishes.
	defer rows.Close()

	// Create an empty slice to hold our transactions.
	var transactions = make([]models.Transaction, 0)

	for rows.Next() {

		var transaction models.Transaction

		err := rows.Scan(
			&transaction.ID,
			&transaction.Title,
			&transaction.Amount,
			&transaction.Category,
			&transaction.Type,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	// Check whether PostgreSQL encountered an error
	// while reading the rows.
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Return all transactions.
	return transactions, nil
}

// AddTransactionToDB saves a new transaction
// into the PostgreSQL database.
func AddTransactionToDB(
	pool *pgxpool.Pool,
	transaction models.Transaction,
) (models.Transaction, error) {

	// Insert the transaction into PostgreSQL.
	// RETURNING id gives us the ID PostgreSQL generated.
	err := pool.QueryRow(
		context.Background(),
		`
		INSERT INTO transactions
			(title, amount, category, type)
		VALUES
			($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
		`,
		transaction.Title,
		transaction.Amount,
		transaction.Category,
		transaction.Type,
	).Scan(&transaction.ID,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)

	// Stop if PostgreSQL returned an error.
	if err != nil {
		return models.Transaction{}, err
	}

	// Return the transaction with its new database ID.
	return transaction, nil
}

// GetTransactionByIDFromDB retrieves one transaction
// from PostgreSQL using its ID.
func GetTransactionByIDFromDB(
	pool *pgxpool.Pool,
	id int,
) (models.Transaction, error) {

	// Create a transaction to receive
	// the data returned by PostgreSQL.
	var transaction models.Transaction

	// Find the transaction with the requested ID.
	err := pool.QueryRow(
		context.Background(),
		`
		SELECT id, title, amount, category, type, created_at, updated_at
		FROM transactions
		WHERE id = $1
		`,
		id,
	).Scan(
		&transaction.ID,
		&transaction.Title,
		&transaction.Amount,
		&transaction.Category,
		&transaction.Type,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)

	// Return an error if the transaction
	// could not be found or another database error occurred.
	if err != nil {
		return models.Transaction{}, err
	}

	// Return the transaction.
	return transaction, nil
}

// UpdateTransactionInDB updates an existing transaction
// in the PostgreSQL database.
func UpdateTransactionInDB(
	pool *pgxpool.Pool,
	id int,
	transaction models.Transaction,
) (models.Transaction, error) {

	// Update the transaction whose ID matches the given ID.
	err := pool.QueryRow(
		context.Background(),
		`
		UPDATE transactions
		SET title = $1,
		    amount = $2,
		    category = $3,
		    type = $4,
			updated_at = NOW()
		WHERE id = $5
		RETURNING id, title, amount, category, type, created_at, updated_at
		`,
		transaction.Title,
		transaction.Amount,
		transaction.Category,
		transaction.Type,
		id,
	).Scan(
		&transaction.ID,
		&transaction.Title,
		&transaction.Amount,
		&transaction.Category,
		&transaction.Type,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)

	// Return an error if the update failed.
	if err != nil {
		return models.Transaction{}, err
	}

	// Return the updated transaction.
	return transaction, nil
}

// DeleteTransactionFromDB removes a transaction
// from the PostgreSQL database.
func DeleteTransactionFromDB(
	pool *pgxpool.Pool,
	id int,
) error {

	// Delete the transaction with the matching ID.
	_, err := pool.Exec(
		context.Background(),
		`
		DELETE FROM transactions
		WHERE id = $1
		`,
		id,
	)

	// Return any database error.
	if err != nil {
		return err
	}

	// Delete was successful.
	return nil
}

// GetFinancialSummaryFromDB calculates total income,
// total expenses, and the current balance from PostgreSQL.
func GetFinancialSummaryFromDB(
	pool *pgxpool.Pool,
) (float64, float64, float64, error) {

	// Variables to receive the totals from PostgreSQL.
	var totalIncome float64
	var totalExpenses float64

	// Calculate income and expenses separately.
	err := pool.QueryRow(
		context.Background(),
		`
		SELECT
			COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0)
		FROM transactions
		`,
	).Scan(
		&totalIncome,
		&totalExpenses,
	)

	// Stop if PostgreSQL returned an error.
	if err != nil {
		return 0, 0, 0, err
	}

	// Calculate the balance in Go.
	balance := totalIncome - totalExpenses

	// Return all three values.
	return totalIncome, totalExpenses, balance, nil
}
