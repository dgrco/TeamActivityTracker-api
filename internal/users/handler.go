package users

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type Handler struct {
	service *Service
}

// Construct a new user Handler
func NewHandler(service *Service) *Handler {
	return &Handler{service}
}

// Sets up routes for users.
func (h *Handler) RegisterRoutes(router *echo.Group) {
	// Create user group
	userRouter := router.Group("/user")

	// Get a user given their ID (set in URL parameter)
	userRouter.GET("/:id", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Test")
	})
}
