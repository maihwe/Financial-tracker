package utils

import (
	"net/mail"
	"strings"
)

// ValidateUserRegistration checks whether the
// registration information is acceptable.
func ValidateUserRegistration(
	email string,
	password string,
) string {

	email = strings.TrimSpace(email)

	if email == "" {
		return "Email is required"
	}

	_, err := mail.ParseAddress(email)

	if err != nil {
		return "Invalid email address"
	}

	if password == "" {
		return "Password is required"
	}

	if len(password) < 8 {
		return "Password must be at least 8 characters"
	}

	return ""
}
