package users

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dgrco/TeamActivityTracker-api/internal/auth"
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
	// Load JWT_SECRET
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		fmt.Fprintf(os.Stderr, "ERROR: JWT_SECRET environment variable is required.\n")
		os.Exit(1)
	}

	// Create user group
	userRouter := router.Group("/user")
	userRouter.Use(auth.JWTMiddleware(jwtSecret))

	// Get a user given their ID (set in URL parameter)
	userRouter.GET("/:id", func(c *echo.Context) error {
		authID, ok := c.Get("user_id").(string)
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
		}

		user, err := h.service.GetUser(c.Request().Context(), authID, c.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
		}
		return c.JSON(http.StatusAccepted, user)
	})
}
