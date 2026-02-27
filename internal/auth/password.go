package auth

import (
	"crypto/sha256"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

// Hashes `password` using the bcrypt algorithm.
// If this hashing fails it will return a non-nil error.
func HashPassword(password string) (string, error) {
	// First compute SHA256 checksum of the password.
	// This ensures that passwords longer than 72 bytes
	// aren't truncated.
	hash := sha256.Sum256([]byte(password))
	// hashStr is not necessary, but useful for debugging and perf loss is negligible.
	// Also, it converts each byte to 2 bytes (hex), so 32 * 2 = 64 (still under the 72 byte max)
	hashStr := hex.EncodeToString(hash[:])

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(hashStr), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

// Checks if a password matches a hashed password.
// If this check fails it will return a non-nil error.
func CheckPassword(password string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
