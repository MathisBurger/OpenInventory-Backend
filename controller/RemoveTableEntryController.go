package controller

import (
	"encoding/json"
	"fmt"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	OwnSQL "github.com/MathisBurger/OpenInventory-Backend/mysql"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

func RemoveTableEntryController(c *fiber.Ctx) error {
	raw := string(c.Body())
	obj := models.RemoveTableEntryRequestModel{}
	err := json.Unmarshal([]byte(raw), &obj)
	if err != nil {
		utils.LogError("[RemoveTableEntryController.go, 17, InputError] " + err.Error())
		res, err := models.GetJSONResponse("Invaild JSON body", "alert alert-danger", "error", "None", 200)
		if err != nil {
			utils.LogError("[RemoveTableEntryController.go, 20, ParsingError] " + err.Error())
		}
		return c.Send(res)
	}
	if !checkRemoveTableEntryRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if OwnSQL.MysqlLoginWithToken(obj.Username, obj.Password, obj.Token) {
		fmt.Println(obj)
		conn := OwnSQL.GetConn()
		stmt, _ := conn.Prepare("SELECT `min-perm-lvl` FROM `inv_tables` WHERE `name`=?;")
		type cacheStruct struct {
			MinPermLvl int `json:"min-perm-lvl"`
		}
		resp, err := stmt.Query(obj.TableName)
		if err != nil {
			utils.LogError("[DeleteTableController.go, 32, SQL-ScanningError] " + err.Error())
		}
		minPermLvl := 0
		for resp.Next() {
			var cache cacheStruct
			err = resp.Scan(&cache.MinPermLvl)
			if err != nil {
				utils.LogError("[DeleteTableController.go, 39, SQL-ScanningError] " + err.Error())
			}
			minPermLvl = cache.MinPermLvl
		}
		defer resp.Close()
		if OwnSQL.CheckUserHasHigherPermission(conn, obj.Username, minPermLvl, "") {
			stmt, _ = conn.Prepare("DELETE FROM `table_" + obj.TableName + "` WHERE `id`=?")
			aff, _ := stmt.Exec(obj.RowID)
			rowsAffected, _ := aff.RowsAffected()
			if rowsAffected == 0 {
				defer resp.Close()
				defer stmt.Close()
				defer conn.Close()
				res, _ := models.GetJSONResponse("EntryID not found", "alert alert-warning", "ok", "None", 200)
				return c.Send(res)
			}
			stmt, _ = conn.Prepare("SELECT `entries` FROM `inv_tables` WHERE `name`=?")
			resp, err := stmt.Query(obj.TableName)
			if err != nil {
				utils.LogError("[RemoveTableEntryController.go, 41, SQL-StatementError] " + err.Error())
			}
			entries := 0
			for resp.Next() {
				var entry OwnSQL.Entries
				err = resp.Scan(&entry.Entries)
				if err != nil {
					utils.LogError("[RemoveTableEntryController.go, 41, SQL-ScanningError] " + err.Error())
				}
				entries = entry.Entries
			}
			entries--
			stmt, _ = conn.Prepare("UPDATE `inv_tables` SET `entries`=? WHERE `name`=?;")
			stmt.Exec(entries, obj.TableName)
			defer resp.Close()
			defer stmt.Close()
			defer conn.Close()
			res, _ := models.GetJSONResponse("Successfully deleted entry", "alert alert-success", "ok", "None", 200)
			return c.Send(res)
		}
		defer stmt.Close()
		defer conn.Close()
		res, _ := models.GetJSONResponse("You do not have the permission perform this", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	res, _ := models.GetJSONResponse("You do not have the permission perform this", "alert alert-danger", "ok", "None", 200)
	return c.Send(res)
}

func checkRemoveTableEntryRequest(obj models.RemoveTableEntryRequestModel) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.TableName != "" && obj.RowID > 0
}
