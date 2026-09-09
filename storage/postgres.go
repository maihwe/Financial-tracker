package storage

import (
	"context"
	"time"

	"financial-tracker/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

// GetAllTransactionsFromDB retrieves all transactions
// belonging to one specific user.
//
// The userID comes from the authenticated session.
// We do not trust the frontend to provide it.
func GetAllTransactionsFromDB(
	pool *pgxpool.Pool,
	userID int,
) ([]models.Transaction, error) {

	rows, err := pool.Query(
		context.Background(),
		`
		SELECT
			id,
			user_id,
			title,
			description,
			amount,
			category_id,
			type,
			transaction_at,
			created_at,
			updated_at
		FROM transactions
		WHERE user_id = $1
		ORDER BY id
		`,
		userID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	transactions := make([]models.Transaction, 0)

	for rows.Next() {

		var transaction models.Transaction

		err := rows.Scan(
			&transaction.ID,
			&transaction.UserID,
			&transaction.Title,
			&transaction.Description,
			&transaction.Amount,
			&transaction.CategoryID,
			&transaction.Type,
			&transaction.TransactionAt,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

// AddTransactionToDB saves a new transaction
// belonging to a specific user.
func AddTransactionToDB(
	pool *pgxpool.Pool,
	transaction models.Transaction,
) (models.Transaction, error) {

	// If the caller did not provide a transaction time,
	// use the current time.
	if transaction.TransactionAt.IsZero() {
		transaction.TransactionAt = time.Now()
	}

	// Insert the transaction into PostgreSQL.
	//
	// UserID comes from the authenticated session,
	// not from the frontend request.
	err := pool.QueryRow(
		context.Background(),
		`
		INSERT INTO transactions
			(
				user_id,
				title,
				description,
				amount,
				category_id,
				type,
				transaction_at
			)
		VALUES
			($1, $2, $3, $4, $5, $6, $7)
		RETURNING
			id,
			user_id,
			title,
			description,
			amount,
			category_id,
			type,
			transaction_at,
			created_at,
			updated_at
		`,
		transaction.UserID,
		transaction.Title,
		transaction.Description,
		transaction.Amount,
		transaction.CategoryID,
		transaction.Type,
		transaction.TransactionAt,
	).Scan(
		&transaction.ID,
		&transaction.UserID,
		&transaction.Title,
		&transaction.Description,
		&transaction.Amount,
		&transaction.CategoryID,
		&transaction.Type,
		&transaction.TransactionAt,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)

	if err != nil {
		return models.Transaction{}, err
	}

	return transaction, nil
}

// GetTransactionByIDFromDB retrieves one transaction
// only when it belongs to the specified user.
func GetTransactionByIDFromDB(
	pool *pgxpool.Pool,
	id int,
	userID int,
) (models.Transaction, error) {

	var transaction models.Transaction

	err := pool.QueryRow(
		context.Background(),
		`
		SELECT
			id,
			user_id,
			title,
			description,
			amount,
			category_id,
			type,
			transaction_at,
			created_at,
			updated_at
		FROM transactions
		WHERE id = $1
		AND user_id = $2
		`,
		id,
		userID,
	).Scan(
		&transaction.ID,
		&transaction.UserID,
		&transaction.Title,
		&transaction.Description,
		&transaction.Amount,
		&transaction.CategoryID,
		&transaction.Type,
		&transaction.TransactionAt,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)

	if err != nil {
		return models.Transaction{}, err
	}

	return transaction, nil
}

// UpdateTransactionInDB updates a transaction
// only when it belongs to the specified user.
func UpdateTransactionInDB(
	pool *pgxpool.Pool,
	id int,
	userID int,
	transaction models.Transaction,
) (models.Transaction, error) {

	err := pool.QueryRow(
		context.Background(),
		`
		UPDATE transactions
		SET
			title = $1,
			description = $2,
			amount = $3,
			category_id = $4,
			type = $5,
			transaction_at = $6,
			updated_at = NOW()
		WHERE id = $7
		AND user_id = $8
		RETURNING
			id,
			user_id,
			title,
			description,
			amount,
			category_id,
			type,
			transaction_at,
			created_at,
			updated_at
		`,
		transaction.Title,
		transaction.Description,
		transaction.Amount,
		transaction.CategoryID,
		transaction.Type,
		transaction.TransactionAt,
		id,
		userID,
	).Scan(
		&transaction.ID,
		&transaction.UserID,
		&transaction.Title,
		&transaction.Description,
		&transaction.Amount,
		&transaction.CategoryID,
		&transaction.Type,
		&transaction.TransactionAt,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)

	if err != nil {
		return models.Transaction{}, err
	}

	return transaction, nil
}

// DeleteTransactionFromDB deletes a transaction
// only when it belongs to the specified user.
func DeleteTransactionFromDB(
	pool *pgxpool.Pool,
	id int,
	userID int,
) error {

	_, err := pool.Exec(
		context.Background(),
		`
		DELETE FROM transactions
		WHERE id = $1
		AND user_id = $2
		`,
		id,
		userID,
	)

	if err != nil {
		return err
	}

	return nil
}

// GetFinancialSummaryFromDB calculates total income,
// total expenses, and balance for one specific user.
func GetFinancialSummaryFromDB(
	pool *pgxpool.Pool,
	userID int,
) (float64, float64, float64, error) {

	var totalIncome float64
	var totalExpenses float64

	err := pool.QueryRow(
		context.Background(),
		`
		SELECT
			COALESCE(
				SUM(
					CASE
						WHEN type = 'income' THEN amount
						ELSE 0
					END
				),
				0
			),
			COALESCE(
				SUM(
					CASE
						WHEN type = 'expense' THEN amount
						ELSE 0
					END
				),
				0
			)
		FROM transactions
		WHERE user_id = $1
		`,
		userID,
	).Scan(
		&totalIncome,
		&totalExpenses,
	)

	if err != nil {
		return 0, 0, 0, err
	}

	balance := totalIncome - totalExpenses

	return totalIncome, totalExpenses, balance, nil
}
