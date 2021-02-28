package controller

import (
	"github.com/MathisBurger/OpenInventory-Backend/config"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

type editTableEntryRequest struct {
	Username  string                 `json:"username"`
	Password  string                 `json:"password"`
	Token     string                 `json:"token"`
	TableName string                 `json:"table_name"`
	ObjectID  int                    `json:"object_id"`
	Row       map[string]interface{} `json:"row"`
}

func EditTableEntryController(c *fiber.Ctx) error {
	obj := new(editTableEntryRequest)
	err := c.BodyParser(obj)
	if err != nil {
		if cfg, _ := config.ParseConfig(); cfg.ServerCFG.LogRequestErrors {
			utils.LogError(err.Error(), "EditTableEntryController.go", 23)
		}
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !checkEditTableEntryRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if actions.MysqlLoginWithToken(obj.Username, obj.Password, obj.Token) {
		table := actions.GetTableByName(obj.TableName)
		conn := actions.GetConn()
		defer conn.Close()
		if actions.CheckUserHasHigherPermission(conn, obj.Username, table.MinPermLvl, "") {
			sql := "UPDATE `table_" + obj.TableName + "` SET "
			first_completed := false
			var values []interface{}
			for k, v := range obj.Row {
				if k != "id" {
					if !first_completed {
						sql += "`" + k + "`=?"
						values = append(values, v)
						first_completed = true
					} else {
						sql += ", `" + k + "`=?"
						values = append(values, v)
					}
				}
			}
			sql += " WHERE `id`=?"
			stmt, err := conn.Prepare(sql)
			defer stmt.Close()
			if err != nil {
				utils.LogError(err.Error(), "EditTableEntryController.go", 56)
			}
			values = append(values, obj.ObjectID)
			_, err = stmt.Exec(values...)
			if err != nil {
				resp, _ := models.GetJSONResponse("Illegal row-map", "alert alert-danger", "ok", "None", 200)
				return c.Send(resp)
			}
			res, _ := models.GetJSONResponse("Successfully updated entry", "alert alert-success", "ok", "None", 200)
			return c.Send(res)
		}
		res, _ := models.GetJSONResponse("You do not have the permission to perform this", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	res, _ := models.GetJSONResponse("You do not have the permission to perform this", "alert alert-danger", "ok", "None", 200)
	return c.Send(res)
}

func checkEditTableEntryRequest(obj *editTableEntryRequest) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.TableName != "" && len(obj.Row) > 0 && obj.ObjectID > 0
}
