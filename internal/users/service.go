package users

type Service struct {
	repository *Repository
}

// Construct a new user Service
func NewService(repository *Repository) *Service {
	return &Service{repository}
}
