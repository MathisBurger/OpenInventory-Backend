package controller

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
	"strconv"
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
		utils.LogError("[RenameTableColumnController.go, 26, InputError] " + err.Error())
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !checkRenameTableColumnRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !actions.MysqlLoginWithToken(obj.Username, obj.Password, obj.Token) {
		res, _ := models.GetJSONResponse("You do not have the permission to perform this", "alert alert-danger", "Failed", "None", 200)
		return c.Send(res)
	}
	conn := actions.GetConn()
	columns := actions.GetTableColumns(obj.Username, obj.Password, obj.Token, obj.TableName)
	if len(columns) == 0 {
		resp, _ := models.GetJSONResponse("You do not have the permission to perform this", "alert alert-danger", "ok", "None", 200)
		return c.Send(resp)
	}
	stmt, _ := conn.Prepare("SELECT `min-perm-lvl` FROM `inv_tables` WHERE `name`=?;")
	type cacheStruct struct {
		MinPermLvl int `json:"min-perm-lvl"`
	}
	resp, err := stmt.Query(obj.TableName)
	if err != nil {
		utils.LogError("[DeleteTableController.go, 50, SQL-ScanningError] " + err.Error())
	}
	minPermLvl := 0
	for resp.Next() {
		var cache cacheStruct
		err = resp.Scan(&cache.MinPermLvl)
		if err != nil {
			utils.LogError("[DeleteTableController.go, 57, SQL-ScanningError] " + err.Error())
		}
		minPermLvl = cache.MinPermLvl
	}
	defer resp.Close()
	if actions.CheckUserHasHigherPermission(conn, obj.Username, minPermLvl, "") {
		for _, val := range columns {
			if val.COLUMN_NAME == obj.OldName {
				var length string
				if val.MAX_LENGTH == nil {
					length = ""
				} else {
					i, _ := val.MAX_LENGTH.(int64)
					length = strconv.Itoa(int(i))
				}
				if val.DATA_TYPE == "int" {
					length = "11"
				}
				stmt, err = conn.Prepare("ALTER TABLE `table_" + obj.TableName + "` CHANGE `" + obj.OldName + "`  `" + obj.NewName + "` " + val.DATA_TYPE +
					"(" + length + ") NULL DEFAULT NULL;")
				if err != nil {
					utils.LogError("[RenameTableColumnController.go, 78, SQL-StatementError] " + err.Error())
					res, _ := models.GetJSONResponse("Error with column name statement", "alert alert-danger", "ok", "None", 200)
					return c.Send(res)
				}
				_, err = stmt.Exec()
				if err != nil {
					utils.LogError("[RenameTableColumnController.go, 84, SQL-StatementError] " + err.Error())
					res, _ := models.GetJSONResponse("Error while changing column name", "alert alert-danger", "ok", "None", 200)
					return c.Send(res)
				}
				defer resp.Close()
				defer stmt.Close()
				defer conn.Close()
				res, _ := models.GetJSONResponse("Successfully changed column name", "alert alert-success", "ok", "None", 200)
				return c.Send(res)
			}
		}
		defer resp.Close()
		defer stmt.Close()
		defer conn.Close()
		res, _ := models.GetJSONResponse("Column not found", "alert alert-warning", "ok", "None", 200)
		return c.Send(res)
	}
	defer resp.Close()
	defer stmt.Close()
	defer conn.Close()
	res, _ := models.GetJSONResponse("You do not have the permission to do this", "alert alert-warning", "ok", "None", 200)
	return c.Send(res)
}

func checkRenameTableColumnRequest(obj RenameTableColumnRequestModel) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.TableName != "" && obj.NewName != "" && obj.OldName != ""
}
