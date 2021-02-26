package controller

import (
	"github.com/MathisBurger/OpenInventory-Backend/config"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

type deleteTableRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Token     string `json:"token"`
	TableName string `json:"table_name"`
}

func DeleteTableController(c *fiber.Ctx) error {
	obj := new(deleteTableRequest)
	err := c.BodyParser(obj)
	if err != nil {
		if cfg, _ := config.ParseConfig(); cfg.ServerCFG.LogRequestErrors {
			utils.LogError(err.Error(), "DeleteTableController.go", 23)
		}
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !checkDeleteTableRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if actions.MysqlLoginWithToken(obj.Username, obj.Password, obj.Token) {
		conn := actions.GetConn()
		defer conn.Close()
		table := actions.GetTableByName(obj.TableName)
		if actions.CheckUserHasHigherPermission(conn, obj.Username, table.MinPermLvl, "") {
			stmt, _ := conn.Prepare("DROP TABLE `table_" + obj.TableName + "`;")
			defer stmt.Close()
			_, err = stmt.Exec()
			if err != nil {
				utils.LogError(err.Error(), "DeleteTableController.go", 39)
				res, _ := models.GetJSONResponse("This table does not exist", "alert alert-warning", "ok", "None", 200)
				return c.Send(res)
			}
			stmt, _ = conn.Prepare("DELETE FROM `inv_tables` WHERE `name`=?")
			defer stmt.Close()
			stmt.Exec(obj.TableName)
			res, _ := models.GetJSONResponse("Successfully deleted table", "alert alert-success", "ok", "None", 200)
			return c.Send(res)
		}
		res, _ := models.GetJSONResponse("You do not have the permission to perform this", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)

	}
	res, _ := models.GetJSONResponse("You do not have the permission to perform this", "alert alert-danger", "ok", "None", 200)
	return c.Send(res)

}

func checkDeleteTableRequest(obj *deleteTableRequest) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.TableName != ""
}
