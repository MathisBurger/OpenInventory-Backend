package controller

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	OwnSQL "github.com/MathisBurger/OpenInventory-Backend/mysql"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

func AddTableEntryController(c *fiber.Ctx) error {
	raw := string(c.Body())
	obj := models.AddTableEntryRequestModel{}
	err := json.Unmarshal([]byte(raw), &obj)
	if err != nil {
		utils.LogError("[AddTableEntryController.go, 16, InputError] " + err.Error())
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !checkAddTableEntryRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	// Login is checked in function below
	status := OwnSQL.AddTableEntry(obj.Username, obj.Password, obj.Token, obj.TableName, obj.Row)
	if status {
		res, _ := models.GetJSONResponse("successful", "alert alert-success", "ok", "None", 200)
		return c.Send(res)
	} else {
		res, _ := models.GetJSONResponse("creation failed", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
}

func checkAddTableEntryRequest(obj models.AddTableEntryRequestModel) bool {
	if obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.TableName != "" && len(obj.Row) > 0 {
		return true
	} else {
		return false
	}
}
