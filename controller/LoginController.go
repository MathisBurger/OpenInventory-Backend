package controller

import (
	"encoding/json"
	"fmt"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	OwnSQL "github.com/MathisBurger/OpenInventory-Backend/mysql"
	"github.com/gofiber/fiber/v2"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginController(c *fiber.Ctx) error {
	raw := string(c.Body())
	obj := LoginRequest{}
	err := json.Unmarshal([]byte(raw), &obj)
	if err != nil {
		response, err := models.GetJsonResponse("Invaild JSON body", "alert alert-danger", "error", "None", 200)
		if err != nil {
			panic(err)
		}
		return c.Send(response)
	}
	status, token := OwnSQL.MySQL_login(obj.Username, obj.Password)
	fmt.Println(token)
	if status {
		response, err := models.GetJsonResponse("Login successful", "alert alert-success", "ok", token, 200)
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
