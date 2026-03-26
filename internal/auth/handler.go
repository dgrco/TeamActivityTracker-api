package auth

import (
	"net/http"

	"github.com/dgrco/TeamActivityTracker-api/internal/environment"
	"github.com/dgrco/TeamActivityTracker-api/internal/errors"
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
			return errors.Respond(
				env,
				c,
				http.StatusInternalServerError,
				"failed to register user",
				err,
			)
		}

		err := h.service.RegisterUser(c.Request().Context(), req)
		if err != nil {
			return errors.Respond(
				env,
				c,
				http.StatusInternalServerError,
				"failed to register user",
				err,
			)
		}
		return c.JSON(http.StatusOK, map[string]string {
			"message": "user saved",
		})
	})

	// Log-in an existing user
	router.POST("/login", func(c *echo.Context) error {
		req := new(LoginRequest)
		if err := c.Bind(req); err != nil {
			return errors.Respond(
				env,
				c,
				http.StatusInternalServerError,
				"failed to log-in user",
				err,
			)
		}

		userID, token, err := h.service.LoginUser(c.Request().Context(), env, req)
		if err != nil {
			return errors.Respond(
				env,
				c,
				http.StatusInternalServerError,
				"failed to log-in user",
				err,
			)
		}

		// Set Authorization header to the generated access token
		// TODO: is this necessary?
		c.Response().Header().Set(echo.HeaderAuthorization, "Bearer "+token)

		// Generate and store new refresh token
		refreshToken := GenerateRefreshToken()
		err = h.service.SaveRefreshToken(c.Request().Context(), userID, refreshToken, DefaultRefreshTokenExpiration())
		if err != nil {
			return errors.Respond(
				env,
				c,
				http.StatusInternalServerError,
				"refresh token error",
				err,
			)
		}

		// Create HTTP-only cookie for refresh token
		cookie := &http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			HttpOnly: true,                 // Prevents JS access
			Secure:   env.CookieSecureMode, // Only can be sent over HTTPS
			SameSite: http.SameSiteLaxMode, // CSRF protection
			Path:     "/",
			MaxAge:   DefaultRefreshTokenDuration,
		}
		c.SetCookie(cookie)

		return c.JSON(http.StatusOK, map[string]string{
			"user_id": userID,
			"access_token": token,
		})
	})

	router.POST("/refresh", func(c *echo.Context) error {
		refreshTokenCookie, err := c.Cookie("refresh_token")
		if err != nil {
			return errors.Respond(
				env,
				c,
				http.StatusBadRequest,
				"refresh token rotation error",
				err,
			)
		}

		refreshToken := refreshTokenCookie.Value
		userID, err := h.service.ValidateRefreshToken(c.Request().Context(), refreshToken)
		if err != nil {
			return errors.Respond(
				env,
				c,
				http.StatusBadRequest,
				"refresh token rotation error",
				err,
			)
		}

		newAccessToken, err := GenerateAccessToken(userID, env.JWTSecret)
		if err != nil {
			return errors.Respond(
				env,
				c,
				http.StatusInternalServerError,
				"refresh token rotation error",
				err,
			)
		}

		c.Response().Header().Set(echo.HeaderAuthorization, newAccessToken)

		return c.JSON(http.StatusOK, map[string]string{
			"access_token": newAccessToken,
		})
	})

	router.POST("/logout", func(c *echo.Context) error {
		refreshTokenCookie, err := c.Cookie("refresh_token")
		if err != nil {
			return errors.Respond(
				env,
				c,
				http.StatusBadRequest,
				"logout failed",
				err,
			)
		}

		refreshToken := refreshTokenCookie.Value
		err = h.service.LogoutUser(c.Request().Context(), refreshToken)
		if err != nil {
			return errors.Respond(
				env,
				c,
				http.StatusInternalServerError,
				"logout failed",
				err,
			)
		}

		return c.JSON(http.StatusOK, map[string]string{
			"message": "user logged out",
		})
	})
}
