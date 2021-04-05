package auth

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/database/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
	"time"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// define lifetime of the session token (long lifetime token)
const sessionLifetime = 24 * 7 * time.Hour

// This endpoint generates an refresh token
// if login credentials are right
func LoginController(c *fiber.Ctx) error {

	data := loginRequest{}
	if err := json.Unmarshal(c.Body(), &data); err != nil {
		return c.SendStatus(400)
	}

	if !checkLoginRequest(data) {
		return c.SendStatus(400)
	}

	if status, _ := actions.MysqlLogin(data.Username, data.Password); !status {
		return c.SendStatus(401)
	}

	// generate token
	tokenStr := utils.Base64(64)
	expires := time.Now().Add(sessionLifetime)
	token := &models.RefreshTokenModel{
		Username: data.Username,
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
	c.Cookie(cookie)

	return c.SendStatus(200)
}

func checkLoginRequest(obj loginRequest) bool {
	return obj.Username != "" && obj.Password != ""
}