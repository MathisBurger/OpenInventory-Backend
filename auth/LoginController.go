package auth

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/accesstoken"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/middleware"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

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

	if _, usr := actions.GetUserByUsername(data.Username); usr.TwoFactor {
		tkn := utils.GenerateToken()
		middleware.TwoFactorCommunicationChannel <- middleware.TwoFactorPair{data.Username, tkn}
		return c.SendString(tkn)
	}

	accesstoken.GenerateRefreshToken(c, data.Username)

	return c.SendStatus(200)
}

func checkLoginRequest(obj loginRequest) bool {
	return obj.Username != "" && obj.Password != ""
}
