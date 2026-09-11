package handlers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"financial-tracker/storage"
	"financial-tracker/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

// fakePasswordResetEmailSender records the email information
// instead of sending a real email.
type fakePasswordResetEmailSender struct {
	to       string
	resetURL string
}

func (fake *fakePasswordResetEmailSender) SendPasswordResetEmail(
	to string,
	resetURL string,
) error {

	fake.to = to
	fake.resetURL = resetURL

	return nil
}

func TestForgotPasswordHandler(t *testing.T) {

	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		t.Skip("DATABASE_URL is not set")
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, databaseURL)

	if err != nil {
		t.Fatalf(
			"failed to create database pool: %v",
			err,
		)
	}

	defer pool.Close()

	email := "forgot-password-handler-test@example.com"

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
		VALUES ($1, $2, $3, 'user')
		RETURNING id
		`,
		"Forgot Password Test User",
		email,
		"test-password-hash",
	).Scan(&userID)

	if err != nil {
		t.Fatalf(
			"failed to create test user: %v",
			err,
		)
	}

	defer func() {
		_, _ = pool.Exec(
			ctx,
			`DELETE FROM users WHERE id = $1`,
			userID,
		)
	}()

	fakeEmailSender := &fakePasswordResetEmailSender{}

	handler := ForgotPasswordHandler(
		pool,
		fakeEmailSender,
	)

	body := `{"email":"` + email + `"}`

	request := httptest.NewRequest(
		http.MethodPost,
		"/forgot-password",
		strings.NewReader(body),
	)

	request.Header.Set(
		"Content-Type",
		"application/json",
	)

	recorder := httptest.NewRecorder()

	handler.ServeHTTP(
		recorder,
		request,
	)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status 200, got %d: %s",
			recorder.Code,
			recorder.Body.String(),
		)
	}

	var response map[string]string

	err = json.NewDecoder(
		recorder.Body,
	).Decode(&response)

	if err != nil {
		t.Fatalf(
			"failed to decode response: %v",
			err,
		)
	}

	expectedMessage :=
		"If the email is registered, a password reset link will be sent."

	if response["message"] != expectedMessage {
		t.Fatalf(
			"unexpected response message: %q",
			response["message"],
		)
	}

	if fakeEmailSender.to != email {
		t.Fatalf(
			"expected reset email to be sent to %q, got %q",
			email,
			fakeEmailSender.to,
		)
	}

	if fakeEmailSender.resetURL == "" {
		t.Fatal(
			"expected password reset URL to be generated",
		)
	}

	var resetToken string
	var expiresAt time.Time

	err = pool.QueryRow(
		ctx,
		`
		SELECT
			password_reset_token,
			password_reset_expires_at
		FROM users
		WHERE id = $1
		`,
		userID,
	).Scan(
		&resetToken,
		&expiresAt,
	)

	if err != nil {
		t.Fatalf(
			"failed to verify reset token: %v",
			err,
		)
	}

	if resetToken == "" {
		t.Fatal(
			"expected reset token to be stored",
		)
	}

	if len(resetToken) != 64 {
		t.Fatalf(
			"expected SHA-256 hash length 64, got %d",
			len(resetToken),
		)
	}

	if !expiresAt.After(time.Now()) {
		t.Fatal(
			"expected reset token to have a future expiration",
		)
	}
}

func TestForgotPasswordHandlerUnknownEmail(t *testing.T) {

	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		t.Skip("DATABASE_URL is not set")
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, databaseURL)

	if err != nil {
		t.Fatalf(
			"failed to create database pool: %v",
			err,
		)
	}

	defer pool.Close()

	fakeEmailSender := &fakePasswordResetEmailSender{}

	handler := ForgotPasswordHandler(
		pool,
		fakeEmailSender,
	)

	body := `{"email":"does-not-exist@example.com"}`

	request := httptest.NewRequest(
		http.MethodPost,
		"/forgot-password",
		strings.NewReader(body),
	)

	request.Header.Set(
		"Content-Type",
		"application/json",
	)

	recorder := httptest.NewRecorder()

	handler.ServeHTTP(
		recorder,
		request,
	)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status 200, got %d",
			recorder.Code,
		)
	}

	var response map[string]string

	err = json.NewDecoder(
		recorder.Body,
	).Decode(&response)

	if err != nil {
		t.Fatalf(
			"failed to decode response: %v",
			err,
		)
	}

	expectedMessage :=
		"If the email is registered, a password reset link will be sent."

	if response["message"] != expectedMessage {
		t.Fatalf(
			"unexpected response message: %q",
			response["message"],
		)
	}

	if fakeEmailSender.to != "" {
		t.Fatal(
			"expected no email to be sent for an unknown address",
		)
	}
}

func TestResetPasswordHandlerInvalidToken(t *testing.T) {

	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		t.Skip("DATABASE_URL is not set")
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, databaseURL)

	if err != nil {
		t.Fatalf(
			"failed to create database pool: %v",
			err,
		)
	}

	defer pool.Close()

	handler := ResetPasswordHandler(pool)

	request := httptest.NewRequest(
		http.MethodPost,
		"/reset-password",
		strings.NewReader(
			"token=invalid-token&password=newpassword123",
		),
	)

	request.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)

	recorder := httptest.NewRecorder()

	handler.ServeHTTP(
		recorder,
		request,
	)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status 400, got %d",
			recorder.Code,
		)
	}
}

func TestResetPasswordHandlerValidToken(t *testing.T) {

	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		t.Skip("DATABASE_URL is not set")
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, databaseURL)

	if err != nil {
		t.Fatalf(
			"failed to create database pool: %v",
			err,
		)
	}

	defer pool.Close()

	email := "reset-password-handler-test@example.com"

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
		VALUES ($1, $2, $3, 'user')
		RETURNING id
		`,
		"Reset Password Test User",
		email,
		"old-password-hash",
	).Scan(&userID)

	if err != nil {
		t.Fatalf(
			"failed to create test user: %v",
			err,
		)
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
		t.Fatalf(
			"failed to generate token: %v",
			err,
		)
	}

	tokenHashBytes :=
		sha256.Sum256([]byte(token))

	tokenHash :=
		hex.EncodeToString(
			tokenHashBytes[:],
		)

	err = storage.SavePasswordResetToken(
		pool,
		userID,
		tokenHash,
		time.Now().Add(30*time.Minute),
	)

	if err != nil {
		t.Fatalf(
			"failed to save reset token: %v",
			err,
		)
	}

	handler := ResetPasswordHandler(pool)

	request := httptest.NewRequest(
		http.MethodPost,
		"/reset-password",
		strings.NewReader(
			"token="+token+"&password=newpassword123",
		),
	)

	request.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)

	recorder := httptest.NewRecorder()

	handler.ServeHTTP(
		recorder,
		request,
	)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status 200, got %d: %s",
			recorder.Code,
			recorder.Body.String(),
		)
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

	if passwordHash == "old-password-hash" {
		t.Fatal(
			"expected password hash to change",
		)
	}

	if resetToken != nil {
		t.Fatal(
			"expected reset token to be invalidated",
		)
	}
}