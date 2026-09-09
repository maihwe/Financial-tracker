package handlers

import (
	"errors"
	"net/http"
	"time"

	"financial-tracker/storage"
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
