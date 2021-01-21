package controller

import (
	"encoding/json"
	"fmt"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	OwnSQL "github.com/MathisBurger/OpenInventory-Backend/mysql"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type RenameTableColumnRequestModel struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Token     string `json:"token"`
	TableName string `json:"table_name"`
	OldName   string `json:"old_name"`
	NewName   string `json:"new_name"`
}

func RenameTableColumnController(c *fiber.Ctx) error {
	raw := string(c.Body())
	obj := RenameTableColumnRequestModel{}
	err := json.Unmarshal([]byte(raw), &obj)
	if err != nil {
		resp, _ := models.GetJsonResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(resp)
	}
	if !checkRenameTableColumnRequest(obj) {
		resp, _ := models.GetJsonResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(resp)
	}
	conn := OwnSQL.GetConn()
	columns := OwnSQL.GetTableColumns(obj.Username, obj.Password, obj.Token, obj.TableName)
	for _, val := range columns {
		if val.COLUMN_NAME == obj.OldName {
			var length string
			if val.MAX_LENGTH == nil {
				length = ""
			} else {
				length = strings.ReplaceAll(strings.Split(fmt.Sprintf("%s", val.MAX_LENGTH), "=")[1], ")", "")
			}
			if val.DATA_TYPE == "int" {
				length = "11"
			}
			stmt, err := conn.Prepare("ALTER TABLE `table_" + obj.TableName + "` CHANGE `" + obj.OldName + "`  `" + obj.NewName + "` " + val.DATA_TYPE +
				"(" + length + ") NULL DEFAULT NULL;")
			if err != nil {
				resp, _ := models.GetJsonResponse("Error with column name statement", "alert alert-danger", "ok", "None", 200)
				return c.Send(resp)
			}
			_, err = stmt.Exec()
			if err != nil {
				resp, _ := models.GetJsonResponse("Error while changing column name", "alert alert-danger", "ok", "None", 200)
				return c.Send(resp)
			}
			resp, _ := models.GetJsonResponse("Successfully changed column name", "alert alert-success", "ok", "None", 200)
			return c.Send(resp)
		}
	}
	resp, _ := models.GetJsonResponse("Column not found", "alert alert-warning", "ok", "None", 200)
	return c.Send(resp)
}

func checkRenameTableColumnRequest(obj RenameTableColumnRequestModel) bool {
	if obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.TableName != "" && obj.NewName != "" && obj.OldName != "" {
		return true
	} else {
		return false
	}
}
