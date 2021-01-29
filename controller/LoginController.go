package controller

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	OwnSQL "github.com/MathisBurger/OpenInventory-Backend/mysql"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

func LoginController(c *fiber.Ctx) error {
	raw := string(c.Body())
	obj := models.LoginRequest{}
	err := json.Unmarshal([]byte(raw), &obj)
	if err != nil {
		utils.LogError("[LoginController.go, 15, InputError] " + err.Error())
		res, err := models.GetJSONResponse("Invaild JSON body", "alert alert-danger", "error", "None", 200)
		if err != nil {
			utils.LogError("[LoginController.go, 19, ParsingError] " + err.Error())
		}
		return c.Send(res)
	}
	if !checkLoginRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	status, token := OwnSQL.MysqlLogin(obj.Username, obj.Password)
	if status {
		res, err := models.GetJSONResponse("Login successful", "alert alert-success", "ok", token, 200)
		if err != nil {
			utils.LogError("[LoginController.go, 31, ParsingError] " + err.Error())
		}
		return c.Send(res)
	}
	res, err := models.GetJSONResponse("Login failed", "alert alert-warning", "ok", "None", 200)
	if err != nil {
		utils.LogError("[LoginController.go, 37, ParsingError] " + err.Error())
	}
	return c.Send(res)

}

func checkLoginRequest(obj models.LoginRequest) bool {
	return obj.Username != "" && obj.Password != ""
}
