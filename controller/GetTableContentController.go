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
		resp, _ := models.GetJsonResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(resp)
	}
	if !checkGetTableContentRequest(obj) {
		resp, _ := models.GetJsonResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(resp)
	}
	if !OwnSQL.MySQL_loginWithToken(obj.Username, obj.Password, obj.Token) {
		resp, _ := models.GetJsonResponse("You do not have the permission to perform this command", "alert alert-warning", "Failed", "None", 200)
		return c.Send(resp)
	} else {
		stmt := "SELECT * FROM `table_" + obj.TableName + "`;"
		conn := OwnSQL.GetConn()
		json, err := utils.QueryToJson(conn, stmt)
		if err != nil {
			utils.LogError("[GetTableContentController.go, 33, SQL-StatementError] " + err.Error())
			resp, _ := models.GetJsonResponse("Invalid table name", "alert alert-danger", "ok", "None", 200)
			return c.Send(resp)
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
