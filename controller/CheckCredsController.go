package controller

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	OwnSQL "github.com/MathisBurger/OpenInventory-Backend/mysql"
	"github.com/gofiber/fiber/v2"
)

type LoginWithTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

func CheckCredsController(c *fiber.Ctx) error {
	raw := string(c.Body())
	obj := LoginWithTokenRequest{}
	err := json.Unmarshal([]byte(raw), &obj)
	if err != nil {
		response, err := models.GetJsonResponse("Invaild JSON body", "alert alert-danger", "error", "None", 200)
		if err != nil {
			panic(err)
		}
		return c.Send(response)
	}
	status := OwnSQL.MySQL_loginWithToken(obj.Username, obj.Password, obj.Token)
	if status {
		response, err := models.GetJsonResponse("Login successful", "alert alert-success", "ok", "None", 200)
		if err != nil {
			panic(err)
		}
		return c.Send(response)
	} else {
		response, err := models.GetJsonResponse("Login failed", "alert alert-warning", "ok", "None", 200)
		if err != nil {
			panic(err)
		}
		return c.Send(response)
	}
}
