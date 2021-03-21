package permission_management

import (
	"encoding/json"
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

////////////////////////////////////////////////////////////////////
//                                                                //
//              RemoveUserFromPermissionGroupController           //
//        This controller removes user from permission group      //
//    It requires removeUserFromPermissionGroupRequest instance   //
//                                                                //
////////////////////////////////////////////////////////////////////
func RemoveUserFromPermissionGroupController(c *fiber.Ctx) error {

	// init and parse the request object
	obj := removeUserFromPermissionGroupRequest{}
	err := json.Unmarshal(c.Body(), &obj)

	// check request
	if err != nil {
		if cfg, _ := config.ParseConfig(); cfg.ServerCFG.LogRequestErrors {
			utils.LogError(err.Error(), "RemoveUserFromPermissionGroupController.go", 25)
		}
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "#d41717", "ok", "None", 200)
		return c.Send(res)
	}
	if !checkRemoveUserFromPermissionGroupRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "#d41717", "ok", "None", 200)
		return c.Send(res)
	}

	// check login
	if !actions.MysqlLoginWithToken(obj.Username, obj.Password, obj.Token) {
		res, _ := models.GetJSONResponse("You do not have the permission to perform this", "#d41717", "Failed", "None", 200)
		return c.Send(res)
	}

	// user should not be allowed to remove default permission
	if obj.PermissionName == "default.everyone" {
		res, _ := models.GetJSONResponse("You can not remove the default permission", "alert alert-warning", "Failed", "None", 200)
		return c.Send(res)
	}

	conn := actions.GetConn()
	defer conn.Close()

	// check permission
	if !actions.CheckUserHasHigherPermission(conn, obj.Username, actions.GetHighestPermission(conn, obj.User), "") {
		res, _ := models.GetJSONResponse("You do not have the permission to perform this", "alert alert-warning", "Failed", "None", 200)
		return c.Send(res)
	}

	_, user := actions.GetUserByUsername(obj.User)
	split := strings.Split(user.Permissions, ";")
	reduced := utils.RemoveValueFromArray(split, obj.PermissionName)

	// calculate new permissions
	newPerms := reduced[0]
	for i, k := range reduced {
		if i == 0 {
			continue
		}
		newPerms += ";" + k
	}
	actions.UpdateUserPermission(obj.User, newPerms)

	res, _ := models.GetJSONResponse("Successfully removed permission from user", "#1db004", "ok", "None", 200)
	return c.Send(res)
}

// checks the request
// struct fields should not be default
func checkRemoveUserFromPermissionGroupRequest(obj removeUserFromPermissionGroupRequest) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.PermissionName != "" && obj.User != ""
}
