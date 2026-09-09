package storage

import (
	"sync"

	"financial-tracker/models"
)

// sessions is our temporary session store.
//
// The map key is the session token.
// The value is the session belonging to that token.
//
// We are using in-memory storage for now.
// Later, we can move sessions to PostgreSQL
// or another persistent session store.
var sessions = make(map[string]models.Session)

// sessionsMutex protects the sessions map.
//
// HTTP requests can happen concurrently, so we
// must prevent two requests from modifying or
// reading the map at the same time.
var sessionsMutex sync.RWMutex

// CreateSession stores a new authenticated session.
func CreateSession(session models.Session) {

	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()

	sessions[session.Token] = session
}

// GetSession retrieves a session using its token.
//
// It returns the session and true when found.
// It returns an empty session and false when
// the token does not exist.
func GetSession(token string) (models.Session, bool) {

	sessionsMutex.RLock()
	defer sessionsMutex.RUnlock()

	session, exists := sessions[token]

	if !exists {
		return models.Session{}, false
	}

	return session, true
}

// DeleteSession removes a session using its token.
func DeleteSession(token string) {

	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()

	delete(sessions, token)
}
