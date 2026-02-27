package auth

import (
	"context"
	"errors"
	"regexp"

	"github.com/dgrco/TeamActivityTracker-api/internal/users"
)

// Constants
const MIN_PASSWORD_LENGTH = 8

// Errors
var ErrEmailNotValid = errors.New("email not valid")
var ErrPasswordTooShort = errors.New("password too short")
var ErrUsernameRequired = errors.New("username is required")

type Service struct {
	userRepo users.Repository
}

func NewService(userRepo users.Repository) *Service {
	return &Service{userRepo: userRepo}
}

// Regex for validation
var emailRegex = regexp.MustCompile(`(?i)^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)

// Registers a new user if all credentials are valid under field constraints
func (s *Service) RegisterUser(ctx context.Context, email string, username string, password string) error {
	emailIsValid := emailRegex.MatchString(email)
	if !emailIsValid {
		return ErrEmailNotValid
	}

	if len(password) < MIN_PASSWORD_LENGTH {
		return ErrPasswordTooShort
	}

	if len(username) == 0 {
		return ErrUsernameRequired
	}

	// Save to users table using `userRepo`
	return s.userRepo.Insert(ctx, email, username, password)
}
