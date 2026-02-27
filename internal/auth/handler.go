package auth

import (
	"fmt"
	"net/http"

	"github.com/dgrco/TeamActivityTracker-api/internal/environment"
	"github.com/labstack/echo/v5"
)

type Handler struct {
	service *Service
}

// Construct a new auth Handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Sets up routes for auth
func (h *Handler) RegisterRoutes(env *environment.Environment, router *echo.Group) {
	// Register a new user
	router.POST("/register", func(c *echo.Context) error {
		req := new(RegisterRequest)
		if err := c.Bind(req); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to bind user credentials: %s", err))
		}

		passwordHash, err := HashPassword(req.Password)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("hashing err: %s", err))
		}

		err = h.service.RegisterUser(c.Request().Context(), req.Email, req.Username, passwordHash)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to register user: %s", err))
		}
		return c.String(http.StatusOK, "user saved")
	})
}
