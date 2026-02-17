package routes

import (
	"github.com/dgrco/TeamActivityTracker-api/internal/handlers"
	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Get("/", handlers.GetRoot)
}
