package controller

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	OwnSQL "github.com/MathisBurger/OpenInventory-Backend/mysql"
	"github.com/gofiber/fiber/v2"
)

type RenameTableRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Token     string `json:"token"`
	TableName string `json:"table_name"`
	NewName   string `json:"new_name"`
}

func RenameTableController(c *fiber.Ctx) error {
	raw := string(c.Body())
	obj := RenameTableRequest{}
	err := json.Unmarshal([]byte(raw), &obj)
	if err != nil {
		resp, _ := models.GetJsonResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(resp)
	}
	if !checkRenameTableRequest(obj) {
		resp, _ := models.GetJsonResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(resp)
	}
	if !OwnSQL.MySQL_loginWithToken(obj.Username, obj.Password, obj.Token) {
		resp, _ := models.GetJsonResponse("You do not have the permission to perform this", "alert alert-danger", "ok", "None", 200)
		return c.Send(resp)
	} else {
		conn := OwnSQL.GetConn()
		stmt, err := conn.Prepare("ALTER TABLE `table_" + obj.TableName + "` RENAME `table_" + obj.NewName + "`;")
		if err != nil {
			panic(err)
		}
		_, err = stmt.Exec()
		if err != nil {
			resp, _ := models.GetJsonResponse("This table does not exists", "alert alert-warning", "ok", "None", 200)
			return c.Send(resp)
		}
		stmt, _ = conn.Prepare("UPDATE `inv_tables` SET `name`=? WHERE `name`=?")
		stmt.Exec(obj.NewName, obj.TableName)
		resp, _ := models.GetJsonResponse("Successfully updated tablename", "alert alert-success", "ok", "None", 200)
		return c.Send(resp)
	}
}

func checkRenameTableRequest(obj RenameTableRequest) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.TableName != "" && obj.NewName != ""
}
