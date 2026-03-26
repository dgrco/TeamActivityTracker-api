package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/dgrco/TeamActivityTracker-api/internal/environment"
	"github.com/dgrco/TeamActivityTracker-api/internal/users"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// ===========================================
// User Repository (start)
// ===========================================
type MockUserRepository struct {
	InsertFunc     func(ctx context.Context, email string, username string, passwordHash string) error
	GetByEmailFunc func(ctx context.Context, email string) (*users.User, string, error)
	GetByIdFunc    func(ctx context.Context, id string) (*users.User, error)

	lastEmail    string
	lastUsername string
	lastPassword string
}

func (mr *MockUserRepository) Insert(ctx context.Context, email string, username string, passwordHash string) error {
	mr.lastEmail = email
	mr.lastUsername = username
	mr.lastPassword = passwordHash

	if mr.InsertFunc != nil {
		return mr.InsertFunc(ctx, email, username, passwordHash)
	}

	return nil
}

func (mr *MockUserRepository) GetByEmail(ctx context.Context, email string) (*users.User, string, error) {
	if mr.GetByEmailFunc != nil {
		return mr.GetByEmailFunc(ctx, email)
	}

	return nil, "", nil
}

func (mr *MockUserRepository) GetById(ctx context.Context, id string) (*users.User, error) {
	if mr.GetByIdFunc != nil {
		return mr.GetByIdFunc(ctx, id)
	}

	return nil, nil
}

// ===========================================
// User Repository (end)
// ===========================================

// ===========================================
// Auth Repository (start)
// ===========================================
type MockAuthRepository struct {
	InsertRefreshTokenHashFunc        func(ctx context.Context, userID string, tokenHash string, expiresAt time.Time) error
	GetUserIDFromRefreshTokenHashFunc func(ctx context.Context, refreshTokenHash string) (string, error)
	RevokeTokenFunc                   func(ctx context.Context, refreshTokenHash string) error

	tokenHash string
}

func (mr *MockAuthRepository) InsertRefreshTokenHash(ctx context.Context, userID string, tokenHash string, expiresAt time.Time) error {
	if mr.InsertRefreshTokenHashFunc != nil {
		return mr.InsertRefreshTokenHashFunc(ctx, userID, tokenHash, expiresAt)
	}

	return nil
}

func (mr *MockAuthRepository) GetUserIDFromRefreshTokenHash(ctx context.Context, refreshTokenHash string) (string, error) {
	if mr.GetUserIDFromRefreshTokenHashFunc != nil {
		return mr.GetUserIDFromRefreshTokenHashFunc(ctx, refreshTokenHash)
	}

	return "", nil
}

func (mr *MockAuthRepository) RevokeToken(ctx context.Context, refreshTokenHash string) error {
	if mr.RevokeTokenFunc != nil {
		return mr.RevokeTokenFunc(ctx, refreshTokenHash)
	}

	return nil
}

// ===========================================
// Auth Repository (end)
// ===========================================

func TestLoginUser(t *testing.T) {
	tests := []struct {
		name       string
		mockFunc   func(ctx context.Context, email string) (*users.User, string, error)
		password   string
		wantErr    error
		wantUserID string
		wantToken  bool
	}{
		{
			name:     "successful login",
			password: "password123",
			mockFunc: func(ctx context.Context, email string) (*users.User, string, error) {
				return &users.User{
					ID:       "user123",
					Email:    "test@example.com",
					Username: "testuser",
				}, "$2a$10$z3XlpKCJHG8JSF9jYd8eMusw7OPOU4yjk.2jWKsU6OmRIWs0kTFtC", nil
			},
			wantErr:    nil,
			wantUserID: "user123",
			wantToken:  true,
		},
		{
			name:     "user not found",
			password: "password123",
			mockFunc: func(ctx context.Context, email string) (*users.User, string, error) {
				return nil, "", pgx.ErrNoRows
			},
			wantErr: pgx.ErrNoRows,
		},
		{
			name:     "invalid password",
			password: "wrongpassword",
			mockFunc: func(ctx context.Context, email string) (*users.User, string, error) {
				return &users.User{
					ID:       "user123",
					Email:    "test@example.com",
					Username: "testuser",
				}, "$2a$10$z3XlpKCJHG8JSF9jYd8eMusw7OPOU4yjk.2jWKsU6OmRIWs0kTFtC", nil
			},
			wantErr: bcrypt.ErrMismatchedHashAndPassword,
		},
		{
			name:     "database error",
			password: "password123",
			mockFunc: func(ctx context.Context, email string) (*users.User, string, error) {
				return nil, "", errors.New("database connection failed")
			},
			wantErr: errors.New("database connection failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := &MockUserRepository{
				GetByEmailFunc: tt.mockFunc,
			}

			service := NewService(nil, mockUserRepo)

			userID, token, err := service.LoginUser(context.Background(), &environment.Environment{
				JWTSecret: "test-secret",
			}, &LoginRequest{
				Email:    "test@example.com",
				Password: tt.password,
			})

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.EqualError(t, tt.wantErr, err.Error())
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantUserID, userID)
			assert.NotNil(t, token)

			if tt.wantToken {
				assert.NotEmpty(t, token)
				assert.Contains(t, token, ".")
			}
		})
	}
}

func TestRegisterUser(t *testing.T) {
	tests := []struct {
		name    string
		req     *RegisterRequest
		wantErr error
	}{
		{
			name: "valid registration",
			req: &RegisterRequest{
				Email:    "test@example.com",
				Username: "testuser",
				Password: "password123",
			},
			wantErr: nil,
		},
		{
			name: "invalid email format",
			req: &RegisterRequest{
				Email:    "invalid-email",
				Username: "testuser",
				Password: "password123",
			},
			wantErr: ErrEmailNotValid,
		},
		{
			name: "password too short",
			req: &RegisterRequest{
				Email:    "test@example.com",
				Username: "testuser",
				Password: "short",
			},
			wantErr: ErrPasswordTooShort,
		},
		{
			name: "empty email",
			req: &RegisterRequest{
				Email:    "",
				Username: "testuser",
				Password: "password123",
			},
			wantErr: ErrEmailNotValid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := &MockUserRepository{}

			service := NewService(nil, mockUserRepo)

			err := service.RegisterUser(context.Background(), tt.req)

			if tt.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, tt.wantErr, err.Error())
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.req.Email, mockUserRepo.lastEmail)
			require.Equal(t, tt.req.Username, mockUserRepo.lastUsername)
			require.NotEqual(t, tt.req.Password, mockUserRepo.lastPassword)
		})
	}
}

func TestSaveRefreshToken(t *testing.T) {
	tests := []struct {
		name         string
		mockFunc     func(ctx context.Context, userID string, tokenHash string, expiresAt time.Time) error
		userID       string
		refreshToken string
		expiresAt    time.Time
		wantErr      error
	}{
		{
			name: "refresh token save success",
			refreshToken: "token",
			userID: "test",
			expiresAt: time.Now().Add(2 * time.Minute),
			mockFunc: func(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error {
				if tokenHash == "token" {
					return ErrHashFailed
				}
				return nil
			},
			wantErr: nil,
		},
		{
			name: "database error",
			userID: "test",
			refreshToken: "token",
			mockFunc: func(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error {
				return errors.New("database error")
			},
			wantErr: errors.New("database error"),
		},
		{
			name: "no user ID provided",
			refreshToken: "token",
			wantErr: ErrNoUser,
		},
		{
			name: "no refresh token",
			userID: "test",
			wantErr: ErrTokenInvalidOrExpired,
		},
		{
			name: "no user or refresh token",
			wantErr: ErrNoUser,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthRepo := &MockAuthRepository{
				InsertRefreshTokenHashFunc: tt.mockFunc,
			}

			service := NewService(mockAuthRepo, nil)

			err := service.SaveRefreshToken(context.Background(), tt.userID, tt.refreshToken, tt.expiresAt)

			if tt.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, tt.wantErr, err.Error())
				return
			}

			require.NoError(t, err)
			require.NotEqual(t, mockAuthRepo.tokenHash, tt.refreshToken)
		})
	}
}
