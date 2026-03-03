package auth

import (
	"crypto/sha256"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

// Hashes `password_sha256` using the bcrypt algorithm.
// If this hashing fails it will return a non-nil error.
// NOTE: Use `GetStringSHA256` for the argument.
func HashPassword(passwordSHA256 string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordSHA256), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

// Performs a SHA256 checksum of a string, and encodes it into a hex string.
func GetStringSHA256(password string) string {
	// First compute SHA256 checksum of the password.
	// This ensures that passwords longer than 72 bytes
	// aren't truncated.
	hash := sha256.Sum256([]byte(password))
	// hashStr is not necessary, but useful for debugging and perf loss is negligible.
	// Also, it converts each byte to 2 bytes (hex), so 32 * 2 = 64 (still under the 72 byte max)
	return hex.EncodeToString(hash[:])
}

// Checks if a password matches a hashed password.
// If this check fails it will return a non-nil error.
func CheckPassword(password string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
