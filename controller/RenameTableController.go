package controller

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	OwnSQL "github.com/MathisBurger/OpenInventory-Backend/mysql"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
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
		utils.LogError("[RenameTableController.go, 24, InputError] " + err.Error())
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
			utils.LogError("[RenameTableController.go, 39, SQL-StatementError] " + err.Error())
		}
		_, err = stmt.Exec()
		if err != nil {
			utils.LogError("[RenameTableController.go, 43, SQL-StatementError] " + err.Error())
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
