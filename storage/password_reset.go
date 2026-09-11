package storage

import (
	"context"
	"errors"
	"time"

	"financial-tracker/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SavePasswordResetToken stores a hashed password-reset token
// and its expiration time for a user.
func SavePasswordResetToken(
	pool *pgxpool.Pool,
	userID int,
	tokenHash string,
	expiresAt time.Time,
) error {

	_, err := pool.Exec(
		context.Background(),
		`
		UPDATE users
		SET
			password_reset_token = $1,
			password_reset_expires_at = $2
		WHERE id = $3
		`,
		tokenHash,
		expiresAt,
		userID,
	)

	return err
}

// GetUserByPasswordResetToken finds a user whose reset token
// matches the supplied token hash and has not expired.
func GetUserByPasswordResetToken(
	pool *pgxpool.Pool,
	tokenHash string,
) (models.User, error) {

	var user models.User

	err := pool.QueryRow(
		context.Background(),
		`
		SELECT
			id,
			name,
			email,
			password_hash,
			role,
			created_at
		FROM users
		WHERE password_reset_token = $1
		  AND password_reset_expires_at > $2
		`,
		tokenHash,
		time.Now(),
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, errors.New(
				"invalid or expired password reset token",
			)
		}

		return models.User{}, err
	}

	return user, nil
}

// ResetPasswordInDB replaces the user's password hash
// and invalidates the password-reset token.
func ResetPasswordInDB(
	pool *pgxpool.Pool,
	userID int,
	passwordHash string,
) error {

	_, err := pool.Exec(
		context.Background(),
		`
		UPDATE users
		SET
			password_hash = $1,
			password_reset_token = NULL,
			password_reset_expires_at = NULL
		WHERE id = $2
		`,
		passwordHash,
		userID,
	)

	return err
}
