package controller

import (
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

type addTableEntryRequest struct {
	Username  string                 `json:"username"`
	Password  string                 `json:"password"`
	Token     string                 `json:"token"`
	TableName string                 `json:"table_name"`
	Row       map[string]interface{} `json:"row"`
}

func AddTableEntryController(c *fiber.Ctx) error {
	obj := new(addTableEntryRequest)
	err := c.BodyParser(obj)
	if err != nil {
		utils.LogError(err.Error(), "AddTableEntryController.go", 22)
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !checkAddTableEntryRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	status := actions.AddTableEntry(obj.Username, obj.Password, obj.Token, obj.TableName, obj.Row)
	if status {
		res, _ := models.GetJSONResponse("successful", "alert alert-success", "ok", "None", 200)
		return c.Send(res)
	}
	res, _ := models.GetJSONResponse("creation failed", "alert alert-danger", "ok", "None", 200)
	return c.Send(res)
}

func checkAddTableEntryRequest(obj *addTableEntryRequest) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.TableName != "" && len(obj.Row) > 0
}
