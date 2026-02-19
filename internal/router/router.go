package router

import "github.com/gofiber/fiber/v3"

// Represents each router based on each version of the API
type VersionedRouters struct {
	V1 fiber.Router
}

// Gets a VersionedRouters object reference that contains
// each versioned router in the API.
func GetVersionedRouters(app *fiber.App) *VersionedRouters {
	routers := &VersionedRouters{}
	apiRouter := app.Group("/api")

	// Create versioned routers
	v1Router := apiRouter.Group("/v1")

	// Connect routers to associated VersionedRouters object
	routers.V1 = v1Router

	return routers
}
