package main

import (
	"github.com/dgrco/TeamActivityTracker-api/internal/db"
	"github.com/dgrco/TeamActivityTracker-api/internal/routes"
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

	// Setup routes
	routes.SetupRoutes(app)

	// Listen to port 3000
	app.Listen(":3000")
}
