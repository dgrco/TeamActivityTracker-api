package router

import (
	"github.com/labstack/echo/v5"
)

// Represents each router based on each version of the API
type VersionedRouters struct {
	V1 *echo.Group
}

// Gets a VersionedRouters object reference that contains
// each versioned router in the API.
func GetVersionedRouters(e *echo.Echo) *VersionedRouters {
	routers := &VersionedRouters{}
	apiRouter := e.Group("/api")

	// Create versioned routers
	v1Router := apiRouter.Group("/v1")

	// Connect routers to associated VersionedRouters object
	routers.V1 = v1Router

	return routers
}
