package main

import (
	"github.com/dgrco/TeamActivityTracker-api/internal/db"
	"github.com/dgrco/TeamActivityTracker-api/internal/router"
	"github.com/dgrco/TeamActivityTracker-api/internal/users"
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

func main() {
	// Setup Fiber
	app := fiber.New()

	// Load .env file if it exists
	_ = godotenv.Load()

	// Setup database connections
	pool := db.SetupDatabase()
	defer pool.Close()

	// Get versioned routers
	routers := router.GetVersionedRouters(app)

	// Wire dependencies of features
	// Repository -> Service -> Handler
	userRepository := users.NewRepository(pool)
	userService := users.NewService(userRepository)
	userHandler := users.NewHandler(userService)
	userHandler.RegisterRoutes(routers.V1)

	// Listen to port 3000
	app.Listen(":3000")
}
