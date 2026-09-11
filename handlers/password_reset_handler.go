package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"os"
	"strings"
	"time"

	"financial-tracker/services"
	"financial-tracker/storage"
	"financial-tracker/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PasswordResetEmailSender defines the behavior required
// to send a password-reset email.
type PasswordResetEmailSender interface {
	SendPasswordResetEmail(to string, resetURL string) error
}

// ForgotPasswordHandler starts the password-reset process.
func ForgotPasswordHandler(
	pool *pgxpool.Pool,
	emailSenders ...PasswordResetEmailSender,
) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set(
			"Content-Type",
			"application/json",
		)

		if r.Method != http.MethodPost {
			http.Error(
				w,
				"Method not allowed",
				http.StatusMethodNotAllowed,
			)
			return
		}

		var request struct {
			Email string `json:"email"`
		}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(
				w,
				"Invalid JSON",
				http.StatusBadRequest,
			)
			return
		}

		email := strings.TrimSpace(request.Email)

		if email == "" {
			http.Error(
				w,
				"Email is required",
				http.StatusBadRequest,
			)
			return
		}

		user, err := storage.GetUserByEmailFromDB(
			pool,
			email,
		)

		// Do not reveal whether the email exists.
		if err != nil {

			if errors.Is(err, pgx.ErrNoRows) {
				json.NewEncoder(w).Encode(
					map[string]string{
						"message": "If the email is registered, a password reset link will be sent.",
					},
				)
				return
			}

			http.Error(
				w,
				"Unable to process request",
				http.StatusInternalServerError,
			)
			return
		}

		token, err := utils.GenerateSessionToken()
		if err != nil {
			http.Error(
				w,
				"Unable to create reset token",
				http.StatusInternalServerError,
			)
			return
		}

		// Store only a SHA-256 hash of the reset token.
		tokenBytes := sha256.Sum256([]byte(token))
		tokenHash := hex.EncodeToString(tokenBytes[:])

		// The reset token is valid for 30 minutes.
		expiresAt := time.Now().Add(30 * time.Minute)

		err = storage.SavePasswordResetToken(
			pool,
			user.ID,
			tokenHash,
			expiresAt,
		)

		if err != nil {
			http.Error(
				w,
				"Unable to process request",
				http.StatusInternalServerError,
			)
			return
		}

		var emailService PasswordResetEmailSender

		if len(emailSenders) > 0 {
			emailService = emailSenders[0]
		} else {
			emailService = services.NewEmailService()
		}

		if emailService == nil {
			http.Error(
				w,
				"Unable to send password reset email",
				http.StatusInternalServerError,
			)
			return
		}

		baseURL := os.Getenv("APP_BASE_URL")

		if baseURL == "" {
			baseURL = "http://" + r.Host
		}

		resetURL :=
			strings.TrimRight(baseURL, "/") +
				"/reset-password?token=" +
				token
		err = emailService.SendPasswordResetEmail(
			user.Email,
			resetURL,
		)

		if err != nil {
			http.Error(
				w,
				"Unable to send password reset email",
				http.StatusInternalServerError,
			)
			return
		}

		json.NewEncoder(w).Encode(
			map[string]string{
				"message": "If the email is registered, a password reset link will be sent.",
			},
		)
	}
}

// ResetPasswordPageHandler displays the password-reset form.
func ResetPasswordPageHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			http.Error(
				w,
				"Method not allowed",
				http.StatusMethodNotAllowed,
			)
			return
		}

		token := strings.TrimSpace(
			r.URL.Query().Get("token"),
		)

		if token == "" {
			http.Error(
				w,
				"Reset token is required",
				http.StatusBadRequest,
			)
			return
		}

		tmpl := template.Must(
			template.New("reset-password").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Reset Password</title>
</head>
<body>
	<h1>Reset Password</h1>

	<form method="POST" action="/reset-password">
		<input
			type="hidden"
			name="token"
			value="{{.Token}}"
		>

		<label for="password">
			New password
		</label>

		<input
			type="password"
			id="password"
			name="password"
			required
			minlength="8"
		>

		<button type="submit">
			Reset Password
		</button>
	</form>
</body>
</html>
`),
		)

		err := tmpl.Execute(
			w,
			map[string]string{
				"Token": token,
			},
		)

		if err != nil {
			http.Error(
				w,
				"Unable to display reset page",
				http.StatusInternalServerError,
			)
			return
		}
	}
}

// ResetPasswordHandler changes a user's password
// using a valid password-reset token.
func ResetPasswordHandler(pool *pgxpool.Pool) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(
				w,
				"Method not allowed",
				http.StatusMethodNotAllowed,
			)
			return
		}

		token := strings.TrimSpace(
			r.FormValue("token"),
		)

		password := r.FormValue("password")

		if token == "" {
			http.Error(
				w,
				"Reset token is required",
				http.StatusBadRequest,
			)
			return
		}

		tokenBytes := sha256.Sum256([]byte(token))
		tokenHash := hex.EncodeToString(tokenBytes[:])

		user, err := storage.GetUserByPasswordResetToken(
			pool,
			tokenHash,
		)

		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusBadRequest,
			)
			return
		}

		validationError := utils.ValidateUserRegistration(
			user.Email,
			password,
		)

		if validationError != "" {
			http.Error(
				w,
				validationError,
				http.StatusBadRequest,
			)
			return
		}

		passwordHash, err := utils.HashPassword(password)
		if err != nil {
			http.Error(
				w,
				"Unable to reset password",
				http.StatusInternalServerError,
			)
			return
		}

		err = storage.ResetPasswordInDB(
			pool,
			user.ID,
			passwordHash,
		)

		if err != nil {
			http.Error(
				w,
				"Unable to reset password",
				http.StatusInternalServerError,
			)
			return
		}

		// Invalidate all existing sessions after
		// a successful password reset.
		storage.DeleteSessionsForUser(user.ID)

		w.Header().Set(
			"Content-Type",
			"application/json",
		)

		json.NewEncoder(w).Encode(
			map[string]string{
				"message": "Password reset successfully",
			},
		)
	}
}

// ResetPasswordRouteHandler handles both GET and POST
// requests for the password-reset route.
func ResetPasswordRouteHandler(pool *pgxpool.Pool) http.HandlerFunc {

	pageHandler := ResetPasswordPageHandler()
	resetHandler := ResetPasswordHandler(pool)

	return func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {

		case http.MethodGet:
			pageHandler(w, r)

		case http.MethodPost:
			resetHandler(w, r)

		default:
			http.Error(
				w,
				"Method not allowed",
				http.StatusMethodNotAllowed,
			)
		}
	}
}
