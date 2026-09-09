package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"financial-tracker/models"
	"financial-tracker/storage"
	"financial-tracker/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// RegistrationRequest represents the data
// required to create a new account.
type RegistrationRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterHandler handles user registration.
func RegisterHandler(pool *pgxpool.Pool) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set(
			"Content-Type",
			"application/json",
		)

		// Registration only accepts POST.
		if r.Method != http.MethodPost {

			http.Error(
				w,
				"Method not allowed",
				http.StatusMethodNotAllowed,
			)

			return
		}

		// Read the registration data.
		var request RegistrationRequest

		err := json.NewDecoder(r.Body).Decode(
			&request,
		)

		if err != nil {

			http.Error(
				w,
				"Invalid JSON",
				http.StatusBadRequest,
			)

			return
		}

		// Validate the email and password.
		validationError :=
			utils.ValidateUserRegistration(
				request.Email,
				request.Password,
			)

		if validationError != "" {

			http.Error(
				w,
				validationError,
				http.StatusBadRequest,
			)

			return
		}

		// Hash the password before it reaches
		// the database.
		passwordHash, err :=
			utils.HashPassword(request.Password)

		if err != nil {

			http.Error(
				w,
				"Could not secure password",
				http.StatusInternalServerError,
			)

			return
		}

		// Build the user model.
		user := models.User{
			Email: strings.ToLower(
				strings.TrimSpace(request.Email),
			),
			PasswordHash: passwordHash,
		}

		// Save the user.
		createdUser, err :=
			storage.CreateUserInDB(
				pool,
				user,
			)

		if err != nil {

			// PostgreSQL error code 23505 means
			// a unique value already exists.
			var pgError *pgconn.PgError

			if errors.As(err, &pgError) &&
				pgError.Code == "23505" {

				http.Error(
					w,
					"Email is already registered",
					http.StatusConflict,
				)

				return
			}

			http.Error(
				w,
				"Could not create account",
				http.StatusInternalServerError,
			)

			return
		}

		// Never send the password hash back
		// to the client.
		createdUser.PasswordHash = ""

		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(
			createdUser,
		)
	}
}

func LoginHandler(pool *pgxpool.Pool) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set(
			"Content-Type",
			"application/json",
		)

		// Login only accepts POST.
		if r.Method != http.MethodPost {

			http.Error(
				w,
				"Method not allowed",
				http.StatusMethodNotAllowed,
			)

			return
		}

		// Read the login data.
		var request RegistrationRequest

		err := json.NewDecoder(r.Body).Decode(
			&request,
		)

		if err != nil {

			http.Error(
				w,
				"Invalid JSON",
				http.StatusBadRequest,
			)

			return
		}

		// Validate the email and password.
		validationError :=
			utils.ValidateUserRegistration(
				request.Email,
				request.Password,
			)

		if validationError != "" {

			http.Error(
				w,
				validationError,
				http.StatusBadRequest,
			)

			return
		}

		// Find the account using the email.
		user, err :=
			storage.GetUserByEmailFromDB(
				pool,
				request.Email,
			)

		if err != nil {

			http.Error(
				w,
				"Invalid email or password",
				http.StatusUnauthorized,
			)

			return
		}

		// Compare the supplied password with
		// the stored bcrypt password hash.
		if !utils.CheckPassword(
			request.Password,
			user.PasswordHash,
		) {

			http.Error(
				w,
				"Invalid email or password",
				http.StatusUnauthorized,
			)

			return
		}

		// Generate a secure random session token.
		token, err :=
			utils.GenerateSessionToken()

		if err != nil {

			http.Error(
				w,
				"Could not create session",
				http.StatusInternalServerError,
			)

			return
		}

		// Create a session for this user.
		session := models.Session{
			Token:     token,
			UserID:    user.ID,
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}

		storage.CreateSession(session)

		// Send the session token to the browser
		// as an HttpOnly cookie.
		http.SetCookie(
			w,
			&http.Cookie{
				Name:     "session_token",
				Value:    token,
				Path:     "/",
				HttpOnly: true,
				Secure:   false,
				SameSite: http.SameSiteLaxMode,
				Expires:  session.ExpiresAt,
			},
		)

		// Never send the password hash to the client.
		user.PasswordHash = ""

		json.NewEncoder(w).Encode(user)
	}
}

// GetUserByIDHandler retrieves a user account
// by ID.
//
// Authentication will be added later so that
// users can only retrieve their own account.
func GetUserByIDHandler(pool *pgxpool.Pool) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// This handler will be completed when
		// authentication and sessions are added.

		_ = pool

		http.Error(
			w,
			"Authentication is required",
			http.StatusUnauthorized,
		)
	}
}

// Keep pgx imported for database error handling
// compatibility with the project's PostgreSQL layer.
var _ = pgx.ErrNoRows
