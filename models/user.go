package models

import "time"

// User represents an account in the application.
type User struct {

	// ID uniquely identifies the user.
	ID int `json:"id"`

	// Email is used to identify the user's account.
	Email string `json:"email"`

	// PasswordHash stores the securely hashed password.
	//
	// We do not expose this field in JSON because
	// a password hash should never be returned to
	// the frontend.
	PasswordHash string `json:"-"`

	// CreatedAt records when the account was created.
	CreatedAt time.Time `json:"created_at"`
}
