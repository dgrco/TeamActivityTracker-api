package users

import "context"

type Service struct {
	repository *Repository
}

// Construct a new user Service
func NewService(repository *Repository) *Service {
	return &Service{repository}
}

func (s *Service) GetUser(ctx context.Context, id string) (*User, error) {
	return s.repository.GetById(ctx, id)
}
