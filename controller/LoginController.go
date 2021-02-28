package controller

import (
	"github.com/MathisBurger/OpenInventory-Backend/config"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

func LoginController(c *fiber.Ctx) error {
	obj := new(models.LoginRequest)
	err := c.BodyParser(obj)
	if err != nil {
		if cfg, _ := config.ParseConfig(); cfg.ServerCFG.LogRequestErrors {
			utils.LogError(err.Error(), "LoginController.go", 16)
		}
		res, _ := models.GetJSONResponse("Invaild JSON body", "alert alert-danger", "error", "None", 200)
		return c.Send(res)
	}
	if !checkLoginRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	status, token := actions.MysqlLogin(obj.Username, obj.Password)
	if status {
		res, _ := models.GetJSONResponse("Login successful", "alert alert-success", "ok", token, 200)
		return c.Send(res)
	}
	res, _ := models.GetJSONResponse("Login failed", "alert alert-warning", "ok", "None", 200)
	return c.Send(res)

}

func checkLoginRequest(obj *models.LoginRequest) bool {
	return obj.Username != "" && obj.Password != ""
}
