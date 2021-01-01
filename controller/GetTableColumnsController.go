package controller

import (
	"encoding/json"
	"fmt"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	OwnSQL "github.com/MathisBurger/OpenInventory-Backend/mysql"
	"github.com/gofiber/fiber/v2"
)

func GetTableColumnsController(c *fiber.Ctx) error {
	raw := string(c.Body())
	obj := models.GetTableColumnsRequestModel{}
	err := json.Unmarshal([]byte(raw), &obj)
	if err != nil {
		resp, _ := models.GetJsonResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(resp)
	}
	columns := OwnSQL.GetTableColumns(obj.Username, obj.Password, obj.Token, obj.TableName)
	if fmt.Sprintf("%T", columns) == "bool" {
		resp, _ := models.GetJsonResponse("Error while fetching Array", "alert alert-danger", "ok", "None", 200)
		return c.Send(resp)
	} else {
		return c.JSON(models.GetTableColumnsResponseModel{
			Message: "successful",
			Alert:   "alert alert-success",
			Columns: columns,
		})
	}
}
