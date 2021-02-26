package controller

import (
	"github.com/MathisBurger/OpenInventory-Backend/config"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type deletePermissionGroupRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Token     string `json:"token"`
	GroupName string `json:"group_name"`
}

func DeletePermissionGroupController(c *fiber.Ctx) error {
	obj := new(deletePermissionGroupRequest)
	err := c.BodyParser(obj)
	if err != nil {
		if cfg, _ := config.ParseConfig(); cfg.ServerCFG.LogRequestErrors {
			utils.LogError(err.Error(), "DeletePermissionGroupController.go", 24)
		}
		res, _ := models.GetJSONResponse("Invaild JSON body", "alert alert-danger", "error", "None", 200)
		return c.Send(res)
	}
	if !checkDeletePermissionGroupRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !actions.MysqlLoginWithToken(obj.Username, obj.Password, obj.Token) {
		res, _ := models.GetJSONResponse("You do not have the permission to perform this", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	conn := actions.GetConn()
	defer conn.Close()
	if actions.CheckUserHasHigherPermission(conn, obj.Username, 0, "permission."+obj.GroupName) {
		user := actions.GetUsersByPermission("%permission." + obj.GroupName + "%")
		for _, val := range user {
			split := strings.Split(val.Permissions, ";")
			reduced := utils.RemoveValueFromArray(split, "permission."+obj.GroupName)
			editedPerms := reduced[0]
			for i, val2 := range reduced {
				if i == 0 {
					continue
				}
				editedPerms += ";" + val2
			}
			actions.UpdateUserPermission(val.Username, editedPerms)
		}
		stmt, _ := conn.Prepare("DELETE FROM `inv_permissions` WHERE `name`=?")
		defer stmt.Close()
		_, err = stmt.Exec("permission." + obj.GroupName)
		if err != nil {
			utils.LogError(err.Error(), "DeletePermissionGroupController.go", 55)
		}
		res, _ := models.GetJSONResponse("Successfully deleted PermissionGroup", "alert alert-success", "ok", "None", 200)
		return c.Send(res)
	}
	res, _ := models.GetJSONResponse("You do not have the permission to perform this", "alert alert-danger", "ok", "None", 200)
	return c.Send(res)
}

func checkDeletePermissionGroupRequest(obj *deletePermissionGroupRequest) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.GroupName != ""
}
