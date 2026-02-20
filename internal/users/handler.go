package users

import "github.com/gofiber/fiber/v3"

type Handler struct {
	service *Service
}

// Construct a new user Handler
func NewHandler(service *Service) *Handler {
	return &Handler{service}
}

// Sets up routes for users.
func (h *Handler) RegisterRoutes(router fiber.Router) {
	// Create user group
	userRouter := router.Group("/user")

	// Get a user given their ID (set in URL parameter)
	userRouter.Get("/:id", func(c fiber.Ctx) error {
		return c.SendString("Hello, " + c.Params("id")) // TODO
	})
}
