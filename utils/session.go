package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateSessionToken creates a cryptographically
// secure random token for an authenticated session.
func GenerateSessionToken() (string, error) {

	// Create 32 random bytes.
	//
	// 32 bytes = 256 bits of randomness.
	bytes := make([]byte, 32)

	_, err := rand.Read(bytes)

	if err != nil {
		return "", err
	}

	// Convert the random bytes into a hexadecimal
	// string that can safely be used as a cookie value.
	return hex.EncodeToString(bytes), nil
}
