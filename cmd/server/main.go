package main

import (
	"net/http"

	"github.com/dgrco/TeamActivityTracker-api/internal/auth"
	"github.com/dgrco/TeamActivityTracker-api/internal/db"
	"github.com/dgrco/TeamActivityTracker-api/internal/environment"
	"github.com/dgrco/TeamActivityTracker-api/internal/router"
	"github.com/dgrco/TeamActivityTracker-api/internal/users"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	// Setup Echo
	e := echo.New()

	// Load environment
	env := environment.Load()

	// CORS Middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{env.WebAppServerURL},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	// Setup database connections
	pool := db.SetupDatabase(env)
	defer pool.Close()

	// Get versioned routers
	routers := router.GetVersionedRouters(e)

	// Create router domain groups
	userRouter := routers.V1.Group("/user", auth.JWTMiddleware(env.JWTSecret))
	authRouter := routers.V1.Group("/auth")

	// Wire dependencies of features

	// User
	userRepository := users.NewPostgresRepository(pool)
	userService := users.NewService(userRepository)
	userHandler := users.NewHandler(userService)
	userHandler.RegisterRoutes(userRouter)

	// Auth
	authRepository := auth.NewPostgresRepository(pool)
	authService := auth.NewService(authRepository, userRepository)
	authHandler := auth.NewHandler(authService)
	authHandler.RegisterRoutes(env, authRouter)

	// Listen to port set by PORT environment variable
	if err := e.Start(":" + env.Port); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
