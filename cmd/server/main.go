package main

import (
	"github.com/dgrco/TeamActivityTracker-api/internal/auth"
	"github.com/dgrco/TeamActivityTracker-api/internal/db"
	"github.com/dgrco/TeamActivityTracker-api/internal/environment"
	"github.com/dgrco/TeamActivityTracker-api/internal/router"
	"github.com/dgrco/TeamActivityTracker-api/internal/users"
	"github.com/labstack/echo/v5"
)

func main() {
	// Setup Echo
	e := echo.New()

	// Load environment
	env := environment.Load()

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
	authService := auth.NewService(userRepository)
	authHandler := auth.NewHandler(authService)
	authHandler.RegisterRoutes(env, authRouter)

	// Listen to port set by PORT environment variable
	if err := e.Start(":" + env.Port); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
