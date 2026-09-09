package storage

import (
	"context"
	"strings"

	"financial-tracker/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

// CreateUserInDB creates a new user account
// in the PostgreSQL database.
//
// The password must already be securely hashed
// before this function is called.
func CreateUserInDB(
	pool *pgxpool.Pool,
	user models.User,
) (models.User, error) {

	// Normalize the email before storing it.
	email := strings.ToLower(
		strings.TrimSpace(user.Email),
	)

	err := pool.QueryRow(
		context.Background(),
		`
		INSERT INTO users
			(
				email,
				password_hash
			)
		VALUES
			($1, $2)
		RETURNING
			id,
			email,
			password_hash,
			created_at
		`,
		email,
		user.PasswordHash,
	).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// GetUserByEmailFromDB finds a user using their email address.
//
// This will later be used during login.
func GetUserByEmailFromDB(
	pool *pgxpool.Pool,
	email string,
) (models.User, error) {

	var user models.User

	// Normalize the email before searching.
	email = strings.ToLower(
		strings.TrimSpace(email),
	)

	err := pool.QueryRow(
		context.Background(),
		`
		SELECT
			id,
			email,
			password_hash,
			created_at
		FROM users
		WHERE email = $1
		`,
		email,
	).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// GetUserByIDFromDB finds a user using their ID.
//
// This will later help us identify the logged-in user
// and display account information.
func GetUserByIDFromDB(
	pool *pgxpool.Pool,
	id int,
) (models.User, error) {

	var user models.User

	err := pool.QueryRow(
		context.Background(),
		`
		SELECT
			id,
			email,
			password_hash,
			created_at
		FROM users
		WHERE id = $1
		`,
		id,
	).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
