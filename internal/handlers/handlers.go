package handlers

import "github.com/gofiber/fiber/v3"

func GetRoot(c fiber.Ctx) error {
	return c.SendString("Hello from Root")
}
