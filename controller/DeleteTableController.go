package controller

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	OwnSQL "github.com/MathisBurger/OpenInventory-Backend/mysql"
	"github.com/gofiber/fiber/v2"
)

func DeleteTableController(c *fiber.Ctx) error {
	raw := string(c.Body())
	obj := models.CreateTableRequestModel{}
	err := json.Unmarshal([]byte(raw), &obj)
	if err != nil {
		resp, _ := models.GetJsonResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(resp)
	}
	if OwnSQL.MySQL_loginWithToken(obj.Username, obj.Password, obj.Token) {
		conn := OwnSQL.GetConn()
		stmt, _ := conn.Prepare("DROP TABLE `table_" + obj.TableName + "`;")
		_, err := stmt.Exec()
		if err != nil {
			resp, _ := models.GetJsonResponse("This table does not exist", "alert alert-warning", "ok", "None", 200)
			return c.Send(resp)
		}
		stmt, _ = conn.Prepare("DELETE FROM `inv_tables` WHERE `name`=?")
		stmt.Exec(obj.TableName)
		resp, _ := models.GetJsonResponse("Successfully deleted table", "alert alert-success", "ok", "None", 200)
		return c.Send(resp)
	} else {
		resp, _ := models.GetJsonResponse("You do not have the permission to perform this", "alert alert-danger", "ok", "None", 200)
		return c.Send(resp)
	}
}
