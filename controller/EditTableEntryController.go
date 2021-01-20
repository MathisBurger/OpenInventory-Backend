package controller

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	OwnSQL "github.com/MathisBurger/OpenInventory-Backend/mysql"
	"github.com/gofiber/fiber/v2"
)

type editTableEntryRequestModel struct {
	Username  string                 `json:"username"`
	Password  string                 `json:"password"`
	Token     string                 `json:"token"`
	TableName string                 `json:"table_name"`
	ObjectID  int                    `json:"object_id"`
	Row       map[string]interface{} `json:"row"`
}

func EditTableEntryController(c *fiber.Ctx) error {
	raw := string(c.Body())
	obj := editTableEntryRequestModel{}
	err := json.Unmarshal([]byte(raw), &obj)
	if err != nil {
		resp, _ := models.GetJsonResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(resp)
	}
	if !checkEditTableEntryRequest(obj) {
		resp, _ := models.GetJsonResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(resp)
	}
	conn := OwnSQL.GetConn()
	sql := "UPDATE `table_" + obj.TableName + "` SET "
	first_completed := false
	var values []interface{}
	for k, v := range obj.Row {
		if !first_completed {
			sql += "`" + k + "`=?"
			values = append(values, v)
			first_completed = true
		} else {
			sql += ", `" + k + "`=?"
			values = append(values, v)
		}
	}
	sql += " WHERE `id`=?"
	stmt, err := conn.Prepare(sql)
	if err != nil {
		panic(err.Error())
	}
	values = append(values, obj.ObjectID)
	_, err = stmt.Exec(values...)
	if err != nil {
		resp, _ := models.GetJsonResponse("Illegal row-map", "alert alert-danger", "ok", "None", 200)
		return c.Send(resp)
	}
	resp, _ := models.GetJsonResponse("Successfully updated entry", "alert alert-success", "ok", "None", 200)
	return c.Send(resp)
}

func checkEditTableEntryRequest(obj editTableEntryRequestModel) bool {
	if obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.TableName != "" && len(obj.Row) > 0 && obj.ObjectID > 0 {
		return true
	} else {
		return false
	}
}
