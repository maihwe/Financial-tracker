package storage

import (
	"context"
	"errors"

	"financial-tracker/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AdminStatistics contains system-wide financial statistics.
type AdminStatistics struct {
	TotalUsers        int     `json:"total_users"`
	TotalTransactions int     `json:"total_transactions"`
	TotalIncome       float64 `json:"total_income"`
	TotalExpenses     float64 `json:"total_expenses"`
	Balance           float64 `json:"balance"`
}

// GetAdminStatisticsFromDB retrieves statistics
// for the entire Financial Tracker system.
func GetAdminStatisticsFromDB(
	pool *pgxpool.Pool,
) (AdminStatistics, error) {

	var statistics AdminStatistics

	err := pool.QueryRow(
		context.Background(),
		`
		SELECT
			(SELECT COUNT(*) FROM users),
			(SELECT COUNT(*) FROM transactions),
			COALESCE(
				(
					SELECT SUM(amount)
					FROM transactions
					WHERE type = 'income'
				),
				0
			),
			COALESCE(
				(
					SELECT SUM(amount)
					FROM transactions
					WHERE type = 'expense'
				),
				0
			)
		`,
	).Scan(
		&statistics.TotalUsers,
		&statistics.TotalTransactions,
		&statistics.TotalIncome,
		&statistics.TotalExpenses,
	)

	if err != nil {
		return AdminStatistics{}, err
	}

	statistics.Balance =
		statistics.TotalIncome -
			statistics.TotalExpenses

	return statistics, nil
}

// UpdateUserRoleInDB changes the role of a user.
func UpdateUserRoleInDB(
	pool *pgxpool.Pool,
	userID int,
	role string,
) (models.User, error) {

	var user models.User

	err := pool.QueryRow(
		context.Background(),
		`
		UPDATE users
		SET role = $1
		WHERE id = $2
		RETURNING id, name, email, role, created_at
		`,
		role,
		userID,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, errors.New("user not found")
		}

		return models.User{}, err
	}

	return user, nil
}
