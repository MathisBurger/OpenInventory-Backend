package controller

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	OwnSQL "github.com/MathisBurger/OpenInventory-Backend/mysql"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

type EditTableMinPermLvlRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Token     string `json:"token"`
	TableName string `json:"table_name"`
	NewLvl    int    `json:"new_lvl"`
}

func EditTableMinPermLvlController(c *fiber.Ctx) error {
	raw := string(c.Body())
	obj := EditTableMinPermLvlRequest{}
	err := json.Unmarshal([]byte(raw), &obj)
	if err != nil {
		utils.LogError("[EditTableEntryController.go, 25, InputError] " + err.Error())
		res, _ := models.GetJsonResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !checkEditTableMinPermLvlRequest(obj) {
		res, _ := models.GetJsonResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !OwnSQL.MySQL_loginWithToken(obj.Username, obj.Password, obj.Token) {
		res, _ := models.GetJsonResponse("You do not have the permission to perform this", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	} else {
		conn := OwnSQL.GetConn()
		stmt, err := conn.Prepare("SELECT `min-perm-lvl` FROM `inv_tables` WHERE `name`=?;")
		if err != nil {
			utils.LogError("[EditTableMinPermLvlController.go, 39, SQL-StatementError] " + err.Error())
		}
		resp, err := stmt.Query(obj.TableName)
		if err != nil {
			utils.LogError("[EditTableMinPermLvlController.go, 43, SQL-StatementError] " + err.Error())
		}
		minPermLvl := 0
		type cacheStruct struct {
			MinPermLvl int `json:"min-perm-lvl"`
		}
		for resp.Next() {
			var cache cacheStruct
			err = resp.Scan(&cache.MinPermLvl)
			if err != nil {
				utils.LogError("[EditTableMinPermLvlController.go, 53, SQL-StatementError] " + err.Error())
			}
			minPermLvl = cache.MinPermLvl
		}
		if !OwnSQL.CheckUserHasHigherPermission(conn, obj.Username, minPermLvl, "") {
			res, _ := models.GetJsonResponse("You do not have the permission to perform this", "alert alert-warning", "ok", "None", 200)
			return c.Send(res)
		}
		if !OwnSQL.CheckUserHasHigherPermission(conn, obj.Username, obj.NewLvl, "") {
			res, _ := models.GetJsonResponse("You can not set teh table level higher then yours", "alert alert-danger", "ok", "None", 200)
			return c.Send(res)
		}
		stmt, err = conn.Prepare("UPDATE `inv_tables` SET `min-perm-lvl`=? WHERE `name`=?;")
		if err != nil {
			utils.LogError("[EditTableMinPermLvlController.go, 67, SQL-StatementError] " + err.Error())
		}
		stmt.Exec(obj.NewLvl, obj.TableName)
		res, _ := models.GetJsonResponse("Successfully updated minimum permission level of table", "alert alert-success", "ok", "None", 200)
		return c.Send(res)
	}
}

func checkEditTableMinPermLvlRequest(obj EditTableMinPermLvlRequest) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.TableName != "" && obj.NewLvl > 0
}