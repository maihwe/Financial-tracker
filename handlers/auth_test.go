package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"financial-tracker/models"
	"financial-tracker/storage"
)

// TestGetAuthenticatedUserID verifies that a valid
// session cookie identifies the correct user.
func TestGetAuthenticatedUserID(t *testing.T) {

	token := "test-session-token"

	session := models.Session{
		Token:     token,
		UserID:    1,
		ExpiresAt: time.Now().Add(time.Hour),
	}

	storage.CreateSession(session)

	request := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	request.AddCookie(
		&http.Cookie{
			Name:  "session_token",
			Value: token,
		},
	)

	userID, err :=
		GetAuthenticatedUserID(request)

	if err != nil {
		t.Fatalf(
			"expected valid session, got error: %v",
			err,
		)
	}

	if userID != 1 {
		t.Fatalf(
			"expected user ID 1, got %d",
			userID,
		)
	}

	// Clean up the test session.
	storage.DeleteSession(token)
}

// TestGetAuthenticatedUserIDWithoutCookie verifies
// that a request without a session cookie is rejected.
func TestGetAuthenticatedUserIDWithoutCookie(t *testing.T) {

	request := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	_, err :=
		GetAuthenticatedUserID(request)

	if err == nil {
		t.Fatal("expected authentication error")
	}
}

// TestGetAuthenticatedUserIDExpiredSession verifies
// that an expired session is rejected.
func TestGetAuthenticatedUserIDExpiredSession(t *testing.T) {

	token := "expired-session-token"

	session := models.Session{
		Token:     token,
		UserID:    1,
		ExpiresAt: time.Now().Add(-time.Hour),
	}

	storage.CreateSession(session)

	request := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	request.AddCookie(
		&http.Cookie{
			Name:  "session_token",
			Value: token,
		},
	)

	_, err :=
		GetAuthenticatedUserID(request)

	if err == nil {
		t.Fatal("expected expired session to be rejected")
	}
}
