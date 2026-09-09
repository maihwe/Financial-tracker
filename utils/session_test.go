package utils

import "testing"

// TestGenerateSessionToken verifies that a session token
// is successfully generated.
func TestGenerateSessionToken(t *testing.T) {

	token, err := GenerateSessionToken()

	if err != nil {
		t.Fatalf(
			"GenerateSessionToken returned an error: %v",
			err,
		)
	}

	if token == "" {
		t.Fatal("expected session token, got empty string")
	}

	// 32 random bytes become 64 hexadecimal characters.
	if len(token) != 64 {
		t.Fatalf(
			"expected token length 64, got %d",
			len(token),
		)
	}
}
