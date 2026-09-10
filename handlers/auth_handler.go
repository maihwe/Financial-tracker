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

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RegistrationRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterHandler creates a new user account.
func RegisterHandler(pool *pgxpool.Pool) http.HandlerFunc {

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

		var request RegistrationRequest

		err := json.NewDecoder(
			r.Body,
		).Decode(&request)

		if err != nil {

			http.Error(
				w,
				"Invalid JSON",
				http.StatusBadRequest,
			)

			return
		}

		request.Name =
			strings.TrimSpace(request.Name)

		if request.Name == "" {

			http.Error(
				w,
				"Name is required",
				http.StatusBadRequest,
			)

			return
		}

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

		passwordHash, err :=
			utils.HashPassword(
				request.Password,
			)

		if err != nil {

			http.Error(
				w,
				"Could not secure password",
				http.StatusInternalServerError,
			)

			return
		}

		user := models.User{
			Name: strings.TrimSpace(request.Name),

			Email: strings.ToLower(
				strings.TrimSpace(request.Email),
			),

			PasswordHash: passwordHash,
		}

		createdUser, err :=
			storage.CreateUserInDB(
				pool,
				user,
			)

		if err != nil {

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

		// Never send the password hash to the browser.
		createdUser.PasswordHash = ""

		w.WriteHeader(
			http.StatusCreated,
		)

		json.NewEncoder(w).Encode(
			createdUser,
		)
	}
}

// LoginHandler authenticates an existing user.
func LoginHandler(pool *pgxpool.Pool) http.HandlerFunc {

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

		var request RegistrationRequest

		err := json.NewDecoder(
			r.Body,
		).Decode(&request)

		if err != nil {

			http.Error(
				w,
				"Invalid JSON",
				http.StatusBadRequest,
			)

			return
		}

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

		session := models.Session{
			Token: token,

			UserID: user.ID,

			ExpiresAt: time.Now().Add(
				24 * time.Hour,
			),
		}

		storage.CreateSession(
			session,
		)

		http.SetCookie(
			w,
			&http.Cookie{
				Name: "session_token",

				Value: token,

				Path: "/",

				HttpOnly: true,

				Secure: false,

				SameSite: http.SameSiteLaxMode,

				Expires: session.ExpiresAt,
			},
		)

		// Never send the password hash to the browser.
		user.PasswordHash = ""

		json.NewEncoder(w).Encode(
			user,
		)
	}
}

// LogoutHandler ends the current user's session.
func LogoutHandler(w http.ResponseWriter, r *http.Request) {

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

	// Look for the session cookie.
	cookie, err :=
		r.Cookie("session_token")

	// If there is no cookie, the user is
	// already logged out.
	if err != nil {

		json.NewEncoder(w).Encode(
			map[string]string{
				"message": "Logged out successfully",
			},
		)

		return
	}

	// Remove the session from our
	// in-memory session storage.
	storage.DeleteSession(
		cookie.Value,
	)

	// Remove the cookie from the browser.
	http.SetCookie(
		w,
		&http.Cookie{
			Name: "session_token",

			Value: "",

			Path: "/",

			HttpOnly: true,

			Secure: false,

			SameSite: http.SameSiteLaxMode,

			MaxAge: -1,

			Expires: time.Unix(1, 0),
		},
	)

	json.NewEncoder(w).Encode(
		map[string]string{
			"message": "Logged out successfully",
		},
	)
}
