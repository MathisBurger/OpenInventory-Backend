package controller

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	OwnSQL "github.com/MathisBurger/OpenInventory-Backend/mysql"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type RemoveUserFromPermissionGroupRequest struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	Token          string `json:"token"`
	PermissionName string `json:"permission_name"`
	User           string `json:"user"`
}

func RemoveUserFromPermissionGroupController(c *fiber.Ctx) error {
	raw := string(c.Body())
	obj := RemoveUserFromPermissionGroupRequest{}
	err := json.Unmarshal([]byte(raw), &obj)
	if err != nil {
		utils.LogError("[GetTableContentController.go, 17, InputError] " + err.Error())
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !checkRemoveUserFromPermissionGroupRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !OwnSQL.MysqlLoginWithToken(obj.Username, obj.Password, obj.Token) {
		res, _ := models.GetJSONResponse("You do not have the permission to perform this", "alert alert-danger", "Failed", "None", 200)
		return c.Send(res)
	}
	conn := OwnSQL.GetConn()
	if !OwnSQL.CheckUserHasHigherPermission(conn, obj.Username, OwnSQL.GetHighestPermission(conn, obj.User), "") {
		defer conn.Close()
		res, _ := models.GetJSONResponse("You do not have the permission to perform this", "alert alert-warning", "Failed", "None", 200)
		return c.Send(res)
	}
	stmt, err := conn.Prepare("SELECT `permissions` FROM `inv_users` WHERE `username`=?;")
	if err != nil {
		utils.LogError("[RemoveUserFromPermissionGroupController.go, 43, SQL-StatementError] " + err.Error())
	}
	resp, err := stmt.Query(obj.User)
	if err != nil {
		utils.LogError("[RemoveUserFromPermissionGroupController.go, 47, SQL-StatementError] " + err.Error())
	}
	type permStruct struct {
		Permissions string `json:"permissions"`
	}
	var permissions string
	for resp.Next() {
		var cache permStruct
		err = resp.Scan(&cache.Permissions)
		if err != nil {
			utils.LogError("[RemoveUserFromPermissionGroupController.go, 57, SQL-StatementError] " + err.Error())
		}
		permissions = cache.Permissions
	}
	split := strings.Split(permissions, ";")
	reduced := utils.RemoveValueFromArray(split, obj.PermissionName)
	newPerms := reduced[0]
	for i, k := range reduced {
		if i == 0 {
			continue
		}
		newPerms += ";" + k
	}
	stmt, err = conn.Prepare("UPDATE `inv_users` SET `permissions`=? WHERE `username`=?")
	if err != nil {
		utils.LogError("[RemoveUserFromPermissionGroupController.go, 57, SQL-StatementError] " + err.Error())
	}
	stmt.Exec(newPerms, obj.User)
	defer resp.Close()
	defer stmt.Close()
	defer conn.Close()
	res, _ := models.GetJSONResponse("Successfully removed permission from user", "alert alert-success", "ok", "None", 200)
	return c.Send(res)
}

func checkRemoveUserFromPermissionGroupRequest(obj RemoveUserFromPermissionGroupRequest) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.PermissionName != "" && obj.User != ""
}
