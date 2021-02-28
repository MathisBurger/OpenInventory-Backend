package controller

import (
	"github.com/MathisBurger/OpenInventory-Backend/config"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type removeUserFromPermissionGroupRequest struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	Token          string `json:"token"`
	PermissionName string `json:"permission_name"`
	User           string `json:"user"`
}

func RemoveUserFromPermissionGroupController(c *fiber.Ctx) error {
	obj := new(removeUserFromPermissionGroupRequest)
	err := c.BodyParser(obj)
	if err != nil {
		if cfg, _ := config.ParseConfig(); cfg.ServerCFG.LogRequestErrors {
			utils.LogError(err.Error(), "RemoveUserFromPermissionGroupController.go", 25)
		}
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !checkRemoveUserFromPermissionGroupRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !actions.MysqlLoginWithToken(obj.Username, obj.Password, obj.Token) {
		res, _ := models.GetJSONResponse("You do not have the permission to perform this", "alert alert-danger", "Failed", "None", 200)
		return c.Send(res)
	}
	if obj.PermissionName == "default.everyone" {
		res, _ := models.GetJSONResponse("You can not remove the default permission", "alert alert-warning", "Failed", "None", 200)
		return c.Send(res)
	}
	conn := actions.GetConn()
	defer conn.Close()
	if !actions.CheckUserHasHigherPermission(conn, obj.Username, actions.GetHighestPermission(conn, obj.User), "") {
		res, _ := models.GetJSONResponse("You do not have the permission to perform this", "alert alert-warning", "Failed", "None", 200)
		return c.Send(res)
	}
	_, user := actions.GetUserByUsername(obj.User)
	split := strings.Split(user.Permissions, ";")
	reduced := utils.RemoveValueFromArray(split, obj.PermissionName)
	newPerms := reduced[0]
	for i, k := range reduced {
		if i == 0 {
			continue
		}
		newPerms += ";" + k
	}
	actions.UpdateUserPermission(obj.User, newPerms)
	res, _ := models.GetJSONResponse("Successfully removed permission from user", "alert alert-success", "ok", "None", 200)
	return c.Send(res)
}

func checkRemoveUserFromPermissionGroupRequest(obj *removeUserFromPermissionGroupRequest) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.PermissionName != "" && obj.User != ""
}
