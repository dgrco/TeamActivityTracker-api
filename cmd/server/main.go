package main

import (
	"github.com/dgrco/TeamActivityTracker-api/internal/db"
	"github.com/dgrco/TeamActivityTracker-api/internal/router"
	"github.com/dgrco/TeamActivityTracker-api/internal/users"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
)

func main() {
	// Setup Echo
	e := echo.New()

	// Load .env file if it exists
	_ = godotenv.Load()

	// Setup database connections
	pool := db.SetupDatabase()
	defer pool.Close()

	// Get versioned routers
	routers := router.GetVersionedRouters(e)

	// Wire dependencies of features
	// Repository -> Service -> Handler
	userRepository := users.NewRepository(pool)
	userService := users.NewService(userRepository)
	userHandler := users.NewHandler(userService)
	userHandler.RegisterRoutes(routers.V1)

	// Listen to port 3000
	if err := e.Start(":3000"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
