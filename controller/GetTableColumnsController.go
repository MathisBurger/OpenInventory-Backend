package controller

import (
	"encoding/json"
	"fmt"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	OwnSQL "github.com/MathisBurger/OpenInventory-Backend/mysql"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

func GetTableColumnsController(c *fiber.Ctx) error {
	raw := string(c.Body())
	obj := models.GetTableColumnsRequestModel{}
	err := json.Unmarshal([]byte(raw), &obj)
	if err != nil {
		utils.LogError("[GetTableColumnsController.go, 17, InputError] " + err.Error())
		res, _ := models.GetJsonResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !checkGetTableColumnsRequest(obj) {
		res, _ := models.GetJsonResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	// login is checked in function below
	columns := OwnSQL.GetTableColumns(obj.Username, obj.Password, obj.Token, obj.TableName)
	if fmt.Sprintf("%T", columns) == "bool" {
		res, _ := models.GetJsonResponse("Error while fetching Array", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	return c.JSON(models.GetTableColumnsResponseModel{
		Message: "successful",
		Alert:   "alert alert-success",
		Columns: columns,
	})
}

func checkGetTableColumnsRequest(obj models.GetTableColumnsRequestModel) bool {
	if obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.TableName != "" {
		return true
	} else {
		return false
	}
}
