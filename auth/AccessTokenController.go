package auth

import (
	"time"

	"github.com/MathisBurger/OpenInventory-Backend/accesstoken"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/gofiber/fiber/v2"
)

type accessTokenModel struct {
	Token    string    `json:"token"`
	Deadline time.Time `json:"deadline"`
}

// lifetime of access token (short life token)
const accessTokenLifetime = 5 * time.Minute

// uninitialized generator for access token
var atgenerator accesstoken.Generator

// This endpoint generates a JWT access token
// based on a refresh token
func AccessTokenController(c *fiber.Ctx) error {

	// initialize access token generator
	atgenerator, _ = accesstoken.NewJWTManager("./certs/private.pem", "")

	refreshToken := c.Cookies("refreshToken", "")

	if refreshToken == "" {

		return c.SendStatus(401)
	}

	exists, model := actions.GetRefreshToken(refreshToken)

	if !exists || time.Now().After(model.Deadline) {

		return c.SendStatus(401)
	}

	accessToken, err := atgenerator.Generate(model.Username, accessTokenLifetime)

	if err != nil {

		return c.SendStatus(500)
	}
	return c.JSON(&accessTokenModel{accessToken, time.Now().Add(accessTokenLifetime)})
}
