package users

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockRepo struct {
	GetByIdFunc func(ctx context.Context, id string) (*User, error)
	InsertFunc func(ctx context.Context, email string, username string, passwordHash string) error
}

func (mr *MockRepo) GetById(ctx context.Context, id string) (*User, error) {
	return mr.GetByIdFunc(ctx, id)
}

func (mr *MockRepo) Insert(ctx context.Context, email string, username string, passwordHash string) error {
	return mr.InsertFunc(ctx, email, username, passwordHash)
}

func TestGetUser(t *testing.T) {
	tests := []struct {
		name        string
		requesterID string
		targetID    string
		mockFunc    func(ctx context.Context, id string) (*User, error)
		wantErr     error
		wantUserID  string
	}{
		{
			name:        "forbidden",
			requesterID: "123",
			targetID:    "234",
			wantErr:     ErrForbidden,
		},
		{
			name: "not found",
			requesterID: "123",
			targetID: "123",
			mockFunc: func(ctx context.Context, id string) (*User, error) {
				return nil, ErrUserNotFound
			},
			wantErr: ErrUserNotFound,
		},
		{
			name: "ok",
			requesterID: "123",
			targetID: "123",
			mockFunc: func(ctx context.Context, id string) (*User, error) {
				return &User{ID: id}, nil
			},
			wantUserID: "123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockRepo{GetByIdFunc: tt.mockFunc}

			service := NewService(mockRepo)

			user, err := service.GetUser(context.Background(), tt.requesterID, tt.targetID)

			if tt.wantErr != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, user)
			assert.Equal(t, tt.wantUserID, user.ID)
		})
	}
}
