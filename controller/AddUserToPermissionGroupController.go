package controller

import (
	"github.com/MathisBurger/OpenInventory-Backend/config"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type addUserToPermissionGroupRequest struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	Token      string `json:"token"`
	Permission string `json:"permission"`
	User       string `json:"user"`
}

func AddUserToPermissionGroupController(c *fiber.Ctx) error {
	obj := new(addUserToPermissionGroupRequest)
	err := c.BodyParser(obj)
	if err != nil {
		if cfg, _ := config.ParseConfig(); cfg.ServerCFG.LogRequestErrors {
			utils.LogError(err.Error(), "AddUserToPermissionGroupController.go", 26)
		}
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
	defer conn.Close()
	if actions.CheckUserHasHigherPermission(conn, obj.Username, 0, obj.Permission) {
		userexists, user := actions.GetUserByUsername(obj.User)
		if !userexists {
			res, _ := models.GetJSONResponse("User does not exist", "alert alert-danger", "ok", "None", 200)
			return c.Send(res)
		}
		if utils.ContainsStr(strings.Split(user.Permissions, ";"), obj.Permission) {
			res, _ := models.GetJSONResponse("The user is already member of this group", "alert alert-warning", "ok", "None", 200)
			return c.Send(res)
		}
		permExists, _ := actions.GetPermissionByName(obj.Permission)
		if !permExists {
			res, _ := models.GetJSONResponse("This permissiongroup does not exist", "alert alert-warning", "ok", "None", 200)
			return c.Send(res)
		}
		finalPermissions := user.Permissions + ";" + obj.Permission
		actions.UpdateUserPermission(obj.User, finalPermissions)
		res, _ := models.GetJSONResponse("User added to permissiongroup", "alert alert-success", "ok", "None", 200)
		return c.Send(res)
	}
	res, _ := models.GetJSONResponse("Your permission-level is too low", "alert alert-warning", "ok", "None", 200)
	return c.Send(res)
}

func checkAddUsertoPermissionGroupRequest(obj *addUserToPermissionGroupRequest) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.Permission != "" && obj.User != ""
}
