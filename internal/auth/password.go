package auth

import "golang.org/x/crypto/bcrypt"

// Hashes `password` using the bcrypt algorithm.
// If this hashing fails it will return a non-nil error.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

// Checks if a password matches a hashed password.
// If this check fails it will return a non-nil error.
func CheckPassword(password string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
