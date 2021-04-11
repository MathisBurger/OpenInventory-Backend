package user_management

import (
	twoFactorAuth "github.com/MathisBurger/OpenInventory-Backend/2fa"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/middleware"
	"github.com/gofiber/fiber/v2"
)

// This controller enables two factor
// authentication for the user, who requested
// this endpoint
func EnableTwoFactorController(c *fiber.Ctx) error {

	if ok, ident := middleware.ValidateAccessToken(c); ok {

		actions.UpdateUser2FA(ident, true)
		secret := twoFactorAuth.GenerateSecret()
		twoFactorAuth.GenerateQR(ident, secret)
		actions.Add2FaSession(secret, ident)

		return c.SendFile("./temp-qr/qr.png", true)
	} else {

		return c.SendStatus(401)
	}
}
