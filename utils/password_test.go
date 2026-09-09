package utils

import "testing"

// TestHashPassword verifies that a password is successfully hashed.
func TestHashPassword(t *testing.T) {

	password := "test-password-123"

	hash, err := HashPassword(password)

	if err != nil {
		t.Fatalf(
			"HashPassword returned an error: %v",
			err,
		)
	}

	if hash == "" {
		t.Fatal("expected password hash, got empty string")
	}

	if hash == password {
		t.Fatal("password should not be stored as plain text")
	}
}

// TestCheckPassword verifies that the correct password
// matches the generated hash.
func TestCheckPassword(t *testing.T) {

	password := "test-password-123"

	hash, err := HashPassword(password)

	if err != nil {
		t.Fatalf(
			"HashPassword returned an error: %v",
			err,
		)
	}

	if !CheckPassword(password, hash) {
		t.Fatal("expected correct password to match")
	}
}

// TestCheckPasswordWrongPassword verifies that an incorrect
// password does not match the stored hash.
func TestCheckPasswordWrongPassword(t *testing.T) {

	password := "test-password-123"
	wrongPassword := "wrong-password"

	hash, err := HashPassword(password)

	if err != nil {
		t.Fatalf(
			"HashPassword returned an error: %v",
			err,
		)
	}

	if CheckPassword(wrongPassword, hash) {
		t.Fatal("expected wrong password not to match")
	}
}
