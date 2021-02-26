package controller

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type AddUserToPermissionGroupRequest struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	Token      string `json:"token"`
	Permission string `json:"permission"`
	User       string `json:"user"`
}

func AddUserToPermissionGroupController(c *fiber.Ctx) error {
	raw := string(c.Body())
	obj := AddUserToPermissionGroupRequest{}
	err := json.Unmarshal([]byte(raw), &obj)
	if err != nil {
		utils.LogError("[AddUserToPermissionGroupController.go, 25, InputError] " + err.Error())
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !checkAddUsertoPermissionGroupRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !actions.MysqlLoginWithToken(obj.Username, obj.Password, obj.Token) {
		res, _ := models.GetJSONResponse("Wrong login credentials", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	conn := actions.GetConn()
	if actions.CheckUserHasHigherPermission(conn, obj.Username, 0, obj.Permission) {
		stmt, err := conn.Prepare("SELECT `permissions` FROM `inv_users` WHERE `username`=?")
		if err != nil {
			utils.LogError("[AddUserToPermissionGroupController.go, 41, SQL-StatementError] " + err.Error())
		}
		type cacheStruct struct {
			Permissions string `json:"permissions"`
		}
		resp, err := stmt.Query(obj.User)
		if err != nil {
			utils.LogError("[AddUserToPermissionGroupController.go, 48, SQL-StatementError] " + err.Error())
		}
		var permissions string
		for resp.Next() {
			var cache cacheStruct
			err = resp.Scan(&cache.Permissions)
			if err != nil {
				utils.LogError("[AddUserToPermissionGroupController.go, 55, SQL-StatementError] " + err.Error())
			}
			permissions = cache.Permissions
		}
		defer resp.Close()
		if utils.ContainsStr(strings.Split(permissions, ";"), obj.Permission) {
			res, _ := models.GetJSONResponse("The user is already member of this group", "alert alert-warning", "ok", "None", 200)
			return c.Send(res)
		}
		stmt, err = conn.Prepare("SELECT * FROM `inv_permissions` WHERE `name`=?;")
		if err != nil {
			utils.LogError("[AddUserToPermissionGroupController.go, 66, SQL-StatementError] " + err.Error())
		}
		resp, err = stmt.Query(obj.Permission)
		if err != nil {
			utils.LogError("[AddUserToPermissionGroupController.go, 70, SQL-StatementError] " + err.Error())
		}
		counter := 0
		for resp.Next() {
			counter++
		}
		if counter == 0 {
			res, _ := models.GetJSONResponse("This permissiongroup does not exist", "alert alert-warning", "ok", "None", 200)
			return c.Send(res)
		}
		finalPermissions := permissions + ";" + obj.Permission
		stmt, err = conn.Prepare("UPDATE `inv_users` SET `permissions`=? WHERE `username`=?;")
		if err != nil {
			utils.LogError("[AddUserToPermissionGroupController.go, 83, SQL-StatementError] " + err.Error())
		}
		_, err = stmt.Exec(finalPermissions, obj.User)
		if err != nil {
			utils.LogError("[AddUserToPermissionGroupController.go, 87, SQL-StatementError] " + err.Error())
		}
		defer stmt.Close()
		defer conn.Close()
		res, _ := models.GetJSONResponse("User added to permissiongroup", "alert alert-success", "ok", "None", 200)
		return c.Send(res)
	}
	defer conn.Close()
	res, _ := models.GetJSONResponse("Your permission-level is too low", "alert alert-warning", "ok", "None", 200)
	return c.Send(res)
}

func checkAddUsertoPermissionGroupRequest(obj AddUserToPermissionGroupRequest) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.Permission != "" && obj.User != ""
}
