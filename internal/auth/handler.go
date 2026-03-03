package auth

import (
	"fmt"
	"net/http"
	"time"

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

// Request structs
type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Sets up routes for auth
func (h *Handler) RegisterRoutes(env *environment.Environment, router *echo.Group) {
	// Register a new user
	router.POST("/register", func(c *echo.Context) error {
		req := new(RegisterRequest)
		if err := c.Bind(req); err != nil {
			return echo.NewHTTPError(
				http.StatusInternalServerError,
				fmt.Sprintf("failed to bind user credentials: %s", err),
			)
		}

		err := h.service.RegisterUser(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(
				http.StatusInternalServerError,
				fmt.Sprintf("failed to register user: %s", err),
			)
		}
		return c.String(http.StatusOK, "user saved")
	})

	// Log-in an existing user
	router.POST("/login", func(c *echo.Context) error {
		req := new(LoginRequest)
		if err := c.Bind(req); err != nil {
			return echo.NewHTTPError(
				http.StatusInternalServerError,
				fmt.Sprintf("failed to bind user credentials: %s", err),
			)
		}

		userID, token, err := h.service.LoginUser(c.Request().Context(), env, req)
		if err != nil {
			return echo.NewHTTPError(
				http.StatusInternalServerError,
				fmt.Sprintf("failed to log-in user: %s", err),
			)
		}

		// Set Authorization header to the generated access token
		c.Response().Header().Set(echo.HeaderAuthorization, "Bearer "+token)

		// Generate and store new refresh token
		refreshToken, err := GenerateRefreshToken()
		if err != nil {
			return echo.NewHTTPError(
				http.StatusInternalServerError,
				fmt.Sprintf("failed to generate refresh token: %s", err),
			)
		}

		refreshTokenDuration := 60 * 60 * 24 * 30
		expiresAt := time.Now().UTC().Add(time.Duration(refreshTokenDuration) * time.Second)

		h.service.SaveRefreshToken(c.Request().Context(), userID, refreshToken, expiresAt)

		// Create HTTP-only cookie for refresh token
		cookie := &http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			HttpOnly: true,                 // Prevents JS access
			Secure:   env.CookieSecureMode, // Only can be sent over HTTPS
			SameSite: http.SameSiteLaxMode, // CSRF protection
			Path:     "/",
			MaxAge:   refreshTokenDuration, // 30 days (in seconds)
		}
		c.SetCookie(cookie)

		return c.JSON(http.StatusOK, map[string]string{
			"access_token": token,
		})
	})
}
