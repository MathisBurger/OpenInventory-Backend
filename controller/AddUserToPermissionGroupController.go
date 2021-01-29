package controller

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	OwnSQL "github.com/MathisBurger/OpenInventory-Backend/mysql"
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
		res, _ := models.GetJsonResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !checkAddUsertoPermissionGroupRequest(obj) {
		res, _ := models.GetJsonResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !OwnSQL.MySQL_loginWithToken(obj.Username, obj.Password, obj.Token) {
		res, _ := models.GetJsonResponse("Wrong login credentials", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	conn := OwnSQL.GetConn()
	if OwnSQL.CheckUserHasHigherPermission(conn, obj.Username, 0, obj.Permission) {
		stmt, err := conn.Prepare("SELECT `permissions` FROM `inv_users` WHERE `username`=?")
		if err != nil {
			utils.LogError("[AddUserToPermissionGroupController.go, 37, SQL-StatementError] " + err.Error())
		}
		type cacheStruct struct {
			Permissions string `json:"permissions"`
		}
		resp, err := stmt.Query(obj.User)
		if err != nil {
			utils.LogError("[AddUserToPermissionGroupController.go, 44, SQL-StatementError] " + err.Error())
		}
		var permissions string
		for resp.Next() {
			var cache cacheStruct
			err = resp.Scan(&cache.Permissions)
			if err != nil {
				utils.LogError("[AddUserToPermissionGroupController.go, 51, SQL-StatementError] " + err.Error())
			}
			permissions = cache.Permissions
		}
		defer resp.Close()
		if utils.ContainsStr(strings.Split(permissions, ";"), obj.Permission) {
			res, _ := models.GetJsonResponse("The user is already member of this group", "alert alert-warning", "ok", "None", 200)
			return c.Send(res)
		} else {
			stmt, err = conn.Prepare("SELECT * FROM `inv_permissions` WHERE `name`=?;")
			if err != nil {
				utils.LogError("[AddUserToPermissionGroupController.go, 63, SQL-StatementError] " + err.Error())
			}
			resp, err = stmt.Query(obj.Permission)
			if err != nil {
				utils.LogError("[AddUserToPermissionGroupController.go, 67, SQL-StatementError] " + err.Error())
			}
			counter := 0
			for resp.Next() {
				counter += 1
			}
			if counter == 0 {
				res, _ := models.GetJsonResponse("This permissiongroup does not exist", "alert alert-warning", "ok", "None", 200)
				return c.Send(res)
			}
			finalPermissions := permissions + ";" + obj.Permission
			stmt, err = conn.Prepare("UPDATE `inv_users` SET `permissions`=? WHERE `username`=?;")
			if err != nil {
				utils.LogError("[AddUserToPermissionGroupController.go, 80, SQL-StatementError] " + err.Error())
			}
			_, err = stmt.Exec(finalPermissions, obj.User)
			if err != nil {
				utils.LogError("[AddUserToPermissionGroupController.go, 84, SQL-StatementError] " + err.Error())
			}
			defer stmt.Close()
			defer conn.Close()
			res, _ := models.GetJsonResponse("User added to permissiongroup", "alert alert-success", "ok", "None", 200)
			return c.Send(res)
		}
	} else {
		defer conn.Close()
		res, _ := models.GetJsonResponse("Your permission-level is too low", "alert alert-warning", "ok", "None", 200)
		return c.Send(res)
	}
}

func checkAddUsertoPermissionGroupRequest(obj AddUserToPermissionGroupRequest) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.Permission != "" && obj.User != ""
}
