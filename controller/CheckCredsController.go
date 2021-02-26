package controller

import (
	"github.com/MathisBurger/OpenInventory-Backend/config"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

func CheckCredsController(c *fiber.Ctx) error {
	obj := new(models.LoginWithTokenRequest)
	err := c.BodyParser(obj)
	if err != nil {
		if cfg, _ := config.ParseConfig(); cfg.ServerCFG.LogRequestErrors {
			utils.LogError(err.Error(), "CheckCredsController.go", 17)
		}
		res, _ := models.GetJSONResponse("Invaild JSON body", "alert alert-danger", "error", "None", 200)

		return c.Send(res)
	}
	if !checkCheckCredsRequestModel(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	status := actions.MysqlLoginWithToken(obj.Username, obj.Password, obj.Token)
	if status {
		res, _ := models.GetJSONResponse("Login successful", "alert alert-success", "ok", "None", 200)

		return c.Send(res)
	}
	res, _ := models.GetJSONResponse("Login failed", "alert alert-warning", "ok", "None", 200)
	return c.Send(res)
}

func checkCheckCredsRequestModel(obj *models.LoginWithTokenRequest) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != ""
}
