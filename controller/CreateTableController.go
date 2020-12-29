package controller

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	OwnSQL "github.com/MathisBurger/OpenInventory-Backend/mysql"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func CreateTableController(c *fiber.Ctx) error {
	raw := string(c.Body())
	obj := models.CreateTableRequestModel{}
	err := json.Unmarshal([]byte(raw), &obj)
	if err != nil {
		resp, _ := models.GetJsonResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(resp)
	}
	status := OwnSQL.CreateTable(obj.Username, obj.Password, obj.Token, obj.TableName, parse(obj.RowConfig))
	if status {
		resp, _ := models.GetJsonResponse("successful", "alert alert-success", "ok", "None", 200)
		return c.Send(resp)
	} else {
		resp, _ := models.GetJsonResponse("creation failed", "alert alert-danger", "ok", "None", 200)
		return c.Send(resp)
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
