package auth

import (
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/dgrco/TeamActivityTracker-api/internal/environment"
	"github.com/dgrco/TeamActivityTracker-api/internal/users"
	"golang.org/x/crypto/bcrypt"
)

// Constants
const MIN_PASSWORD_LENGTH = 8

// Errors
var ErrEmailNotValid = errors.New("email not valid")
var ErrPasswordTooShort = errors.New("password too short")
var ErrUsernameRequired = errors.New("username is required")
var ErrHashFailed = errors.New("password hashing failed")
var ErrNoToken = errors.New("no token")
var ErrTokenInvalidOrExpired = errors.New("token invalid or expired")
var ErrNoUser = errors.New("no user")

type Service struct {
	authRepo Repository
	userRepo users.Repository
}

func NewService(authRepo Repository, userRepo users.Repository) *Service {
	return &Service{
		authRepo: authRepo,
		userRepo: userRepo,
	}
}

const DefaultRefreshTokenDuration int = 60 * 60 * 24 * 30 // 30 days (in seconds)

// Return the default expiration time for refresh tokens: 30 days from creation
func DefaultRefreshTokenExpiration() time.Time {
	return time.Now().Add(time.Duration(DefaultRefreshTokenDuration) * time.Second)
}

// Email validation regular expression
var emailRegex = regexp.MustCompile(`(?i)^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)

// Registers a new user if all credentials are valid under field constraints
func (s *Service) RegisterUser(ctx context.Context, req *RegisterRequest) error {
	emailIsValid := emailRegex.MatchString(req.Email)
	if !emailIsValid {
		return ErrEmailNotValid
	}

	if len(req.Password) < MIN_PASSWORD_LENGTH {
		return ErrPasswordTooShort
	}

	if len(req.Email) == 0 {
		return ErrUsernameRequired
	}

	// Hash password
	passwordSHA256 := GetStringSHA256(req.Password)
	hashedPassword, err := HashPassword(passwordSHA256)
	if err != nil {
		return ErrHashFailed
	}

	// Save to users table using `userRepo`
	return s.userRepo.Insert(ctx, req.Email, req.Username, hashedPassword)
}

// Logs-in a user given their email and password.
// Returns the user's ID and an access token (string) on success, and an error on failure.
func (s *Service) LoginUser(ctx context.Context, env *environment.Environment, req *LoginRequest) (string, string, error) {
	user, passwordHash, err := s.userRepo.GetByEmail(ctx, req.Email)

	if err != nil {
		return "", "", err
	}

	passwordSHA256 := GetStringSHA256(req.Password)
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(passwordSHA256))
	if err != nil {
		return "", "", err
	}

	token, err := GenerateAccessToken(user.ID, env.JWTSecret)
	return user.ID, token, err
}

// Save a refresh token to the database
func (s *Service) SaveRefreshToken(ctx context.Context, userID string, refreshToken string, expiresAt time.Time) error {
	if userID == "" {
		return ErrNoUser
	}
	if refreshToken == "" {
		return ErrTokenInvalidOrExpired
	}

	tokenHash := GetStringSHA256(refreshToken)
	return s.authRepo.InsertRefreshTokenHash(ctx, userID, tokenHash, expiresAt)
}

// Validate a refresh token. Rotates the refresh token on success.
// Returns ID of user who owns the refresh token, if valid.
// Returns an error if invalid (empty, revoked, expired, or not found) or if a DB issue occurs.
func (s *Service) ValidateRefreshToken(ctx context.Context, refreshToken string) (string, error) {
	if refreshToken == "" {
		return "", ErrNoToken
	}

	refreshTokenHash := GetStringSHA256(refreshToken)
	uid, err := s.authRepo.GetUserIDFromRefreshTokenHash(ctx, refreshTokenHash)

	if err == nil {
		// Rotate refresh token on success
		newRefreshToken := GenerateRefreshToken()
		err = s.SaveRefreshToken(ctx, uid, newRefreshToken, DefaultRefreshTokenExpiration())
		if err != nil {
			return "", err
		}

		err = s.authRepo.RevokeToken(ctx, refreshTokenHash)
		if err != nil {
			return "", err
		}
	}

	return uid, err
}

// Log out a user.
func (s *Service) LogoutUser(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return ErrNoToken
	}

	refreshTokenHash := GetStringSHA256(refreshToken)
	return s.authRepo.RevokeToken(ctx, refreshTokenHash)
}
