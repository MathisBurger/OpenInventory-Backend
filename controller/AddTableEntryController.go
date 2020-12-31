package controller

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	OwnSQL "github.com/MathisBurger/OpenInventory-Backend/mysql"
	"github.com/gofiber/fiber/v2"
)

func AddTableEntryController(c *fiber.Ctx) error {
	raw := string(c.Body())
	obj := models.AddTableEntryRequestModel{}
	err := json.Unmarshal([]byte(raw), &obj)
	if err != nil {
		resp, _ := models.GetJsonResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(resp)
	}
	status := OwnSQL.AddTableEntry(obj.Username, obj.Password, obj.Token, obj.TableName, obj.Row)
	if status {
		resp, _ := models.GetJsonResponse("successful", "alert alert-success", "ok", "None", 200)
		return c.Send(resp)
	} else {
		resp, _ := models.GetJsonResponse("creation failed", "alert alert-danger", "ok", "None", 200)
		return c.Send(resp)
	}
}
