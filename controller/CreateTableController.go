package controller

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	OwnSQL "github.com/MathisBurger/OpenInventory-Backend/mysql"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
	"strings"
)

// create table endpoint
func CreateTableController(c *fiber.Ctx) error {
	raw := string(c.Body())
	obj := models.CreateTableRequestModel{}
	err := json.Unmarshal([]byte(raw), &obj)
	if err != nil {
		utils.LogError("[CreateTableController.go, 17, InputError] " + err.Error())
		res, _ := models.GetJsonResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !checkCreateTableRequestModel(obj) {
		res, _ := models.GetJsonResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !checkTableNameLength(obj.TableName) {
		res, _ := models.GetJsonResponse("Table name is too long", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	// login creds are checked in function below
	status := OwnSQL.CreateTable(obj.Username, obj.Password, obj.Token, obj.TableName, parse(obj.RowConfig), obj.MinPermLvl)
	if status {
		res, _ := models.GetJsonResponse("successful", "alert alert-success", "ok", "None", 200)
		return c.Send(res)
	} else {
		res, _ := models.GetJsonResponse("creation failed", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
}

func parse(val string) (ans []models.RowConfigModel) {
	arr := strings.Split(val, ",")
	for _, el := range arr {
		raws := strings.ReplaceAll(el, "(", "")
		raws = strings.ReplaceAll(raws, ")", "")
		raws = strings.ReplaceAll(raws, " ", "")
		raws = strings.ReplaceAll(raws, "[", "")
		raws = strings.ReplaceAll(raws, "]", "")
		raws = strings.ReplaceAll(raws, "-", "_")
		spl := strings.Split(raws, ";")
		ans = append(ans, models.RowConfigModel{spl[0], CheckTableName(spl[1])})
	}
	return
}

func CheckTableName(name string) string {
	return strings.ReplaceAll(name, "-", "_")
}

func checkCreateTableRequestModel(obj models.CreateTableRequestModel) bool {
	if obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.TableName != "" && obj.RowConfig != "" {
		return true
	} else {
		return true
	}
}

func checkTableNameLength(name string) bool {
	split := strings.Split(name, "")
	return len(split) < 16 && len(split) != 0
}
