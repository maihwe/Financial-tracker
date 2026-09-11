package storage

import (
	"context"
	"os"
	"testing"
	"time"

	"financial-tracker/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPasswordResetTokenLifecycle(t *testing.T) {
	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		t.Skip("DATABASE_URL is not set")
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatalf("failed to create database pool: %v", err)
	}
	defer pool.Close()

	testEmail := "password-reset-test@example.com"

	var userID int

	err = pool.QueryRow(
		ctx,
		`
		INSERT INTO users (
			name,
			email,
			password_hash,
			role
		)
		VALUES (
			$1,
			$2,
			$3,
			'user'
		)
		RETURNING id
		`,
		"Password Reset Test User",
		testEmail,
		"old-password-hash",
	).Scan(&userID)

	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	defer func() {
		_, _ = pool.Exec(
			ctx,
			`DELETE FROM users WHERE id = $1`,
			userID,
		)
	}()

	token, err := utils.GenerateSessionToken()
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	expiresAt := time.Now().Add(30 * time.Minute)

	err = SavePasswordResetToken(
		pool,
		userID,
		token,
		expiresAt,
	)
	if err != nil {
		t.Fatalf("failed to save reset token: %v", err)
	}

	user, err := GetUserByPasswordResetToken(
		pool,
		token,
	)
	if err != nil {
		t.Fatalf("failed to retrieve valid reset token: %v", err)
	}

	if user.ID != userID {
		t.Fatalf(
			"expected user ID %d, got %d",
			userID,
			user.ID,
		)
	}

	expiredToken := "expired-test-token"

	err = SavePasswordResetToken(
		pool,
		userID,
		expiredToken,
		time.Now().Add(-1*time.Minute),
	)
	if err != nil {
		t.Fatalf("failed to save expired token: %v", err)
	}

	_, err = GetUserByPasswordResetToken(
		pool,
		expiredToken,
	)

	if err == nil {
		t.Fatal("expected expired token to be rejected")
	}

	newPasswordHash := "new-password-hash"

	err = ResetPasswordInDB(
		pool,
		userID,
		newPasswordHash,
	)
	if err != nil {
		t.Fatalf("failed to reset password: %v", err)
	}

	var passwordHash string
	var resetToken *string

	err = pool.QueryRow(
		ctx,
		`
		SELECT
			password_hash,
			password_reset_token
		FROM users
		WHERE id = $1
		`,
		userID,
	).Scan(
		&passwordHash,
		&resetToken,
	)

	if err != nil {
		t.Fatalf(
			"failed to verify password reset: %v",
			err,
		)
	}

	if passwordHash != newPasswordHash {
		t.Fatalf(
			"expected password hash %q, got %q",
			newPasswordHash,
			passwordHash,
		)
	}

	if resetToken != nil {
		t.Fatal("expected reset token to be invalidated")
	}
}
