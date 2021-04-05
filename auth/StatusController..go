package auth

import (
	"github.com/MathisBurger/OpenInventory-Backend/middleware"
	"github.com/gofiber/fiber/v2"
)

// This is an endpoint to test if
// the JWTs are working
func StatusController(c *fiber.Ctx) error {

	if status, ident := middleware.ValidateAccessToken(c); status {

		return c.SendString(ident)
	} else {

		return c.SendStatus(401)
	}
}