package auth

import (
	"fmt"
	twoFactorAuth "github.com/MathisBurger/OpenInventory-Backend/2fa"
	"github.com/MathisBurger/OpenInventory-Backend/accesstoken"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/middleware"
	"github.com/gofiber/fiber/v2"
)

// This controller authorizes the given
// code and checks if it is valid.
// If it is, it returns a valid refreshToken
func TwoFactorAuthController(c *fiber.Ctx) error {

	token := c.Query("token", "")
	username := c.Query("username", "")
	code := c.Query("code", "")

	if token != "" && username != "" && code != "" {

		found := false

		for _, el := range middleware.TwoFactorPairs {
			if username == el.Username && token == el.Token {
				found = true
				break
			}
		}

		if found {

			sessions := actions.GetAll2FaSessionsOfUser(username)
			session := sessions[len(sessions)-1]

			fmt.Println("sess", session)

			if twoFactorAuth.Authenticate(session.Secret, code) {

				accesstoken.GenerateRefreshToken(c, username)

				return c.SendStatus(200)
			} else {

				return c.SendStatus(401)
			}
		} else {

			return c.SendStatus(401)
		}
	} else {

		return c.SendStatus(401)
	}
}
