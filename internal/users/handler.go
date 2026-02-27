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
	// Get a user given their ID (set in URL parameter)
	router.GET("/:id", func(c *echo.Context) error {
		authID, ok := c.Get("user_id").(string)
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
		}

		user, err := h.service.GetUser(c.Request().Context(), authID, c.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
		}

		return c.JSON(http.StatusOK, user)
	})
}
