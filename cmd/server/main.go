package main

import (
	"github.com/dgrco/TeamActivityTracker-api/internal/routes"
	"github.com/gofiber/fiber/v3"
)

func main() {
	app := fiber.New()

	routes.SetupRoutes(app)

	app.Listen(":3000")
}
