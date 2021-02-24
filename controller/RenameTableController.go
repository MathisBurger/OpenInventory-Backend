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
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !checkRenameTableRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !OwnSQL.MysqlLoginWithToken(obj.Username, obj.Password, obj.Token) {
		res, _ := models.GetJSONResponse("You do not have the permission to perform this", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	conn := OwnSQL.GetConn()
	stmt, _ := conn.Prepare("SELECT `min-perm-lvl` FROM `inv_tables` WHERE `name`=?;")
	type cacheStruct struct {
		MinPermLvl int `json:"min-perm-lvl"`
	}
	resp, err := stmt.Query(obj.TableName)
	if err != nil {
		utils.LogError("[DeleteTableController.go, 43, SQL-ScanningError] " + err.Error())
	}
	minPermLvl := 0
	for resp.Next() {
		var cache cacheStruct
		err = resp.Scan(&cache.MinPermLvl)
		if err != nil {
			utils.LogError("[DeleteTableController.go, 50, SQL-ScanningError] " + err.Error())
		}
		minPermLvl = cache.MinPermLvl
	}
	defer resp.Close()
	if OwnSQL.CheckUserHasHigherPermission(conn, obj.Username, minPermLvl, "") {
		stmt, err = conn.Prepare("ALTER TABLE `table_" + obj.TableName + "` RENAME `table_" + obj.NewName + "`;")
		if err != nil {
			utils.LogError("[RenameTableController.go, 58, SQL-StatementError] " + err.Error())
		}
		_, err = stmt.Exec()
		if err != nil {
			utils.LogError("[RenameTableController.go, 62, SQL-StatementError] " + err.Error())
			res, _ := models.GetJSONResponse("This table does not exists", "alert alert-warning", "ok", "None", 200)
			return c.Send(res)
		}
		stmt, _ = conn.Prepare("UPDATE `inv_tables` SET `name`=? WHERE `name`=?")
		stmt.Exec(obj.NewName, obj.TableName)
		defer resp.Close()
		defer stmt.Close()
		defer conn.Close()
		res, _ := models.GetJSONResponse("Successfully updated tablename", "alert alert-success", "ok", "None", 200)
		return c.Send(res)
	}
	defer resp.Close()
	defer stmt.Close()
	defer conn.Close()
	res, _ := models.GetJSONResponse("You do not have the permission to perform this", "alert alert-danger", "ok", "None", 200)
	return c.Send(res)
}

func checkRenameTableRequest(obj RenameTableRequest) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.TableName != "" && obj.NewName != ""
}
