package errors

import (
	"log"

	"github.com/dgrco/TeamActivityTracker-api/internal/environment"
	"github.com/labstack/echo/v5"
)

// Log a detailed error to the server and send the client a generic `message`.
// Logs are formatted as: [ERROR] <method> <url>: <detailed error>
func Respond(env *environment.Environment, c *echo.Context, status int, message string, err error) error {
	if err != nil {
		log.Printf("[ERROR] %s %s: %v\n", c.Request().Method, c.Request().URL.Path, err)
	}

	return echo.NewHTTPError(status, message)
}
