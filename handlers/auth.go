package handlers

import (
	"errors"
	"net/http"
	"time"

	"financial-tracker/models"
	"financial-tracker/storage"

	"github.com/jackc/pgx/v5/pgxpool"
)

// GetAuthenticatedUserID identifies the user associated
// with the session cookie on an HTTP request.
//
// It returns the user's ID when the session is valid.
// It returns an error when the user is not authenticated.
func GetAuthenticatedUserID(r *http.Request) (int, error) {

	// Read the session cookie from the request.
	cookie, err := r.Cookie("session_token")

	if err != nil {
		return 0, errors.New("authentication required")
	}

	// Find the session using the token.
	session, exists :=
		storage.GetSession(cookie.Value)

	if !exists {
		return 0, errors.New("invalid session")
	}

	// Check whether the session has expired.
	if time.Now().After(session.ExpiresAt) {

		// Remove the expired session from storage.
		storage.DeleteSession(session.Token)

		return 0, errors.New("session expired")
	}

	// The session is valid.
	return session.UserID, nil
}

// GetAuthenticatedUser returns the full user associated
// with the current session.
func GetAuthenticatedUser(
	pool *pgxpool.Pool,
	r *http.Request,
) (models.User, error) {

	userID, err :=
		GetAuthenticatedUserID(r)

	if err != nil {
		return models.User{}, err
	}

	user, err :=
		storage.GetUserByIDFromDB(
			pool,
			userID,
		)

	if err != nil {
		return models.User{}, errors.New("user not found")
	}

	return user, nil
}

// RequireAdmin checks whether the current user
// has administrator privileges.
//
// It returns true when the user is an admin
// or super_admin. It writes the appropriate
// HTTP error and returns false when access is denied.
func RequireAdmin(
	pool *pgxpool.Pool,
	w http.ResponseWriter,
	r *http.Request,
) bool {

	user, err :=
		GetAuthenticatedUser(
			pool,
			r,
		)

	if err != nil {

		http.Error(
			w,
			"Authentication required",
			http.StatusUnauthorized,
		)

		return false
	}

	if user.Role != "admin" &&
		user.Role != "super_admin" {

		http.Error(
			w,
			"Admin access required",
			http.StatusForbidden,
		)

		return false
	}

	return true
}
