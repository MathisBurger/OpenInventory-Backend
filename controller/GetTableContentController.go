package controller

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	OwnSQL "github.com/MathisBurger/OpenInventory-Backend/mysql"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func GetTableContentController(c *fiber.Ctx) error {
	raw := string(c.Body())
	obj := models.GetTableContentRequestModel{}
	err := json.Unmarshal([]byte(raw), &obj)
	if err != nil {
		utils.LogError("[GetTableContentController.go, 17, InputError] " + err.Error())
		res, _ := models.GetJsonResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !checkGetTableContentRequest(obj) {
		res, _ := models.GetJsonResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !OwnSQL.MySQL_loginWithToken(obj.Username, obj.Password, obj.Token) {
		res, _ := models.GetJsonResponse("You do not have the permission to perform this", "alert alert-warning", "Failed", "None", 200)
		return c.Send(res)
	} else {
		stmt := "SELECT * FROM `table_" + obj.TableName + "`;"
		conn := OwnSQL.GetConn()
		json, err := utils.QueryToJson(conn, stmt)
		if err != nil {
			utils.LogError("[GetTableContentController.go, 33, SQL-StatementError] " + err.Error())
			res, _ := models.GetJsonResponse("Invalid table name", "alert alert-danger", "ok", "None", 200)
			return c.Send(res)
		}
		defer conn.Close()
		return c.JSON(models.GetTableContentResponseModel{
			Message:    "successful",
			Alert:      "alert alert-success",
			Status:     "ok",
			HttpStatus: 200,
			Elements:   strings.ReplaceAll(strings.ReplaceAll(string(json), "\n", ""), "\t", ""),
		})
	}
}

func checkGetTableContentRequest(obj models.GetTableContentRequestModel) bool {
	if obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.TableName != "" {
		return true
	} else {
		return false
	}
}
