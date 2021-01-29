package controller

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	OwnSQL "github.com/MathisBurger/OpenInventory-Backend/mysql"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

func CheckCredsController(c *fiber.Ctx) error {
	raw := string(c.Body())
	obj := models.LoginWithTokenRequest{}
	err := json.Unmarshal([]byte(raw), &obj)
	if err != nil {
		res, err := models.GetJSONResponse("Invaild JSON body", "alert alert-danger", "error", "None", 200)
		if err != nil {
			utils.LogError("[CheckCredsController.go, 18, InputError] " + err.Error())
		}
		return c.Send(res)
	}
	if !checkCheckCredsRequestModel(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	status := OwnSQL.MysqlLoginWithToken(obj.Username, obj.Password, obj.Token)
	if status {
		res, err := models.GetJSONResponse("Login successful", "alert alert-success", "ok", "None", 200)
		if err != nil {
			utils.LogError("[CheckCredsController.go, 30, ParsingError] " + err.Error())
		}
		return c.Send(res)
	} else {
		res, err := models.GetJSONResponse("Login failed", "alert alert-warning", "ok", "None", 200)
		if err != nil {
			utils.LogError("[CheckCredsController.go, 36, ParsingError] " + err.Error())
		}
		return c.Send(res)
	}
}

func checkCheckCredsRequestModel(obj models.LoginWithTokenRequest) bool {
	if obj.Username != "" && obj.Password != "" && obj.Token != "" {
		return true
	} else {
		return false
	}
}
