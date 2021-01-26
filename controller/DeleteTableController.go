package controller

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	OwnSQL "github.com/MathisBurger/OpenInventory-Backend/mysql"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

func DeleteTableController(c *fiber.Ctx) error {
	raw := string(c.Body())
	obj := models.DeleteTableRequestModel{}
	err := json.Unmarshal([]byte(raw), &obj)
	if err != nil {
		utils.LogError("[DeleteTableController.go, 16, InputError] " + err.Error())
		res, _ := models.GetJsonResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !checkDeleteTableRequest(obj) {
		res, _ := models.GetJsonResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if OwnSQL.MySQL_loginWithToken(obj.Username, obj.Password, obj.Token) {
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
			stmt, _ = conn.Prepare("DROP TABLE `table_" + obj.TableName + "`;")
			_, err = stmt.Exec()
			if err != nil {
				utils.LogError("[DeleteTableController.go, 48, SQL-StatementExecutionError] " + err.Error())
				resp, _ := models.GetJsonResponse("This table does not exist", "alert alert-warning", "ok", "None", 200)
				return c.Send(resp)
			}
			stmt, _ = conn.Prepare("DELETE FROM `inv_tables` WHERE `name`=?")
			stmt.Exec(obj.TableName)
			res, _ := models.GetJsonResponse("Successfully deleted table", "alert alert-success", "ok", "None", 200)
			defer stmt.Close()
			defer conn.Close()
			return c.Send(res)
		} else {
			res, _ := models.GetJsonResponse("You do not have the permission to perform this", "alert alert-danger", "ok", "None", 200)
			return c.Send(res)
		}
	} else {
		res, _ := models.GetJsonResponse("You do not have the permission to perform this", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
}

func checkDeleteTableRequest(obj models.DeleteTableRequestModel) bool {
	if obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.TableName != "" {
		return true
	} else {
		return false
	}
}
