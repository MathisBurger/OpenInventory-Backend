package accesstoken

import (
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/database/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
	"time"
)

// define lifetime of the session token (long lifetime token)
const sessionLifetime = 24 * 7 * time.Hour

// This function generates a refresh token
// based on the username and the fiber
// context
func GenerateRefreshToken(c *fiber.Ctx, username string) {
	// generate token
	tokenStr := utils.Base64(64)
	expires := time.Now().Add(sessionLifetime)
	token := &models.RefreshTokenModel{
		Username: username,
		Token:    tokenStr,
		Deadline: expires,
	}

	actions.AddRefreshToken(token)

	// define cookie
	cookie := new(fiber.Cookie)
	cookie.Name = "refreshToken"
	cookie.Value = tokenStr
	cookie.Expires = expires
	cookie.Secure = false
	cookie.HTTPOnly = true
	cookie.SameSite = "None" // Only for development
	c.Cookie(cookie)
}
