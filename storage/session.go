package storage

import (
	"sync"

	"financial-tracker/models"
)

var sessions = make(map[string]models.Session)
var sessionsMutex sync.RWMutex

func CreateSession(session models.Session) {
	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()

	sessions[session.Token] = session
}

func GetSession(token string) (models.Session, bool) {
	sessionsMutex.RLock()
	defer sessionsMutex.RUnlock()

	session, exists := sessions[token]

	if !exists {
		return models.Session{}, false
	}

	return session, true
}

func DeleteSession(token string) {
	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()

	delete(sessions, token)
}

// DeleteSessionsForUser removes all active sessions
// belonging to a specific user.
func DeleteSessionsForUser(userID int) {

	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()

	for token, session := range sessions {

		if session.UserID == userID {
			delete(sessions, token)
		}
	}
}
