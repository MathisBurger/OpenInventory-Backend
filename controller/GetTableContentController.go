package controller

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/models"
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
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !checkGetTableContentRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !actions.MysqlLoginWithToken(obj.Username, obj.Password, obj.Token) {
		res, _ := models.GetJSONResponse("You do not have the permission to perform this", "alert alert-warning", "Failed", "None", 200)
		return c.Send(res)
	}
	stmt := "SELECT * FROM `table_" + obj.TableName + "`;"
	conn := actions.GetConn()
	json, err := utils.QueryToJson(conn, stmt)
	if err != nil {
		utils.LogError("[GetTableContentController.go, 33, SQL-StatementError] " + err.Error())
		res, _ := models.GetJSONResponse("Invalid table name", "alert alert-danger", "ok", "None", 200)
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

func checkGetTableContentRequest(obj models.GetTableContentRequestModel) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.TableName != ""
}
