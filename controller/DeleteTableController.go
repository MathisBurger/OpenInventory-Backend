package controller

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions/utils"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/gofiber/fiber/v2"
)

func DeleteTableController(c *fiber.Ctx) error {
	raw := string(c.Body())
	obj := models.DeleteTableRequestModel{}
	err := json.Unmarshal([]byte(raw), &obj)
	if err != nil {
		utils.LogError("[DeleteTableController.go, 16, InputError] " + err.Error())
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !checkDeleteTableRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if actions.MysqlLoginWithToken(obj.Username, obj.Password, obj.Token) {
		conn := actions.GetConn()
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
		if actions.CheckUserHasHigherPermission(conn, obj.Username, minPermLvl, "") {
			stmt, _ = conn.Prepare("DROP TABLE `table_" + obj.TableName + "`;")
			_, err = stmt.Exec()
			if err != nil {
				utils.LogError("[DeleteTableController.go, 48, SQL-StatementExecutionError] " + err.Error())
				res, _ := models.GetJSONResponse("This table does not exist", "alert alert-warning", "ok", "None", 200)
				return c.Send(res)
			}
			stmt, _ = conn.Prepare("DELETE FROM `inv_tables` WHERE `name`=?")
			stmt.Exec(obj.TableName)
			res, _ := models.GetJSONResponse("Successfully deleted table", "alert alert-success", "ok", "None", 200)
			defer stmt.Close()
			defer conn.Close()
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

func checkDeleteTableRequest(obj models.DeleteTableRequestModel) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.TableName != ""
}
