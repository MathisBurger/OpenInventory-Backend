package controller

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
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
		utils.LogError("[EditTableEntryController.go, 25, InputError] " + err.Error())
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !checkEditTableEntryRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	conn := actions.GetConn()
	if actions.MysqlLoginWithToken(obj.Username, obj.Password, obj.Token) {
		stmt, _ := conn.Prepare("SELECT `min-perm-lvl` FROM `inv_tables` WHERE `name`=?;")
		type cacheStruct struct {
			MinPermLvl int `json:"min-perm-lvl"`
		}
		resp, err := stmt.Query(obj.TableName)
		if err != nil {
			utils.LogError("[EditTableEntryController.go, 41, SQL-ScanningError] " + err.Error())
		}
		minPermLvl := 0
		for resp.Next() {
			var cache cacheStruct
			err = resp.Scan(&cache.MinPermLvl)
			if err != nil {
				utils.LogError("[EditTableEntryController.go, 48, SQL-ScanningError] " + err.Error())
			}
			minPermLvl = cache.MinPermLvl
		}
		defer resp.Close()
		if actions.CheckUserHasHigherPermission(conn, obj.Username, minPermLvl, "") {
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
			stmt, err = conn.Prepare(sql)
			if err != nil {
				utils.LogError("[EditTableEntryController.go, 72, SQL-StatementError] " + err.Error())
			}
			values = append(values, obj.ObjectID)
			_, err = stmt.Exec(values...)
			if err != nil {
				utils.LogError("[EditTableEntryController.go, 77, SQL-StatementError] " + err.Error())
				resp, _ := models.GetJSONResponse("Illegal row-map", "alert alert-danger", "ok", "None", 200)
				return c.Send(resp)
			}
			defer stmt.Close()
			defer conn.Close()
			res, _ := models.GetJSONResponse("Successfully updated entry", "alert alert-success", "ok", "None", 200)
			return c.Send(res)
		}
		defer stmt.Close()
		defer conn.Close()
		res, _ := models.GetJSONResponse("You do not have the permission to perform this", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	res, _ := models.GetJSONResponse("You do not have the permission to perform this", "alert alert-danger", "ok", "None", 200)
	return c.Send(res)
}

func checkEditTableEntryRequest(obj editTableEntryRequestModel) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.TableName != "" && len(obj.Row) > 0 && obj.ObjectID > 0
}
