package controller

import (
	"github.com/MathisBurger/OpenInventory-Backend/config"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

type removeTableEntryRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Token     string `json:"token"`
	TableName string `json:"table_name"`
	RowID     int    `json:"row_id"`
}

func RemoveTableEntryController(c *fiber.Ctx) error {
	obj := new(removeTableEntryRequest)
	err := c.BodyParser(obj)
	if err != nil {
		if cfg, _ := config.ParseConfig(); cfg.ServerCFG.LogRequestErrors {
			utils.LogError(err.Error(), "RemoveTableEntryController.go", 25)
		}
		res, _ := models.GetJSONResponse("Invaild JSON body", "alert alert-danger", "error", "None", 200)
		return c.Send(res)
	}
	if !checkRemoveTableEntryRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if actions.MysqlLoginWithToken(obj.Username, obj.Password, obj.Token) {
		table := actions.GetTableByName(obj.TableName)
		conn := actions.GetConn()
		defer conn.Close()
		if actions.CheckUserHasHigherPermission(conn, obj.Username, table.MinPermLvl, "") {
			if !actions.DeleteTableEntry(obj.RowID, obj.TableName) {
				res, _ := models.GetJSONResponse("EntryID does not exist", "alert alert-warning", "ok", "None", 200)
				return c.Send(res)
			}
			actions.ChangeNumOfEntrysBy(obj.TableName, -1)
			res, _ := models.GetJSONResponse("Successfully deleted entry", "alert alert-success", "ok", "None", 200)
			return c.Send(res)
		}
		res, _ := models.GetJSONResponse("You do not have the permission perform this", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	res, _ := models.GetJSONResponse("You do not have the permission perform this", "alert alert-danger", "ok", "None", 200)
	return c.Send(res)
}

func checkRemoveTableEntryRequest(obj *removeTableEntryRequest) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.TableName != "" && obj.RowID > 0
}
