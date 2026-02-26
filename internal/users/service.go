package users

import (
	"context"
	"errors"
)

// Errors
var ErrForbidden = errors.New("forbidden")
var ErrUserNotFound = errors.New("user not found")

// Service
type Service struct {
	repository Repository
}

// Construct a new user Service
func NewService(repository Repository) *Service {
	return &Service{repository}
}

// Return a User if the user_id exists in the database and if
// the target ID (set via params) is the same as the authenticated
// user's ID.
func (s *Service) GetUser(ctx context.Context, actorID string, targetID string) (*User, error) {
	if actorID != targetID {
		return nil, ErrForbidden
	}

	return s.repository.GetById(ctx, targetID)
}
