package controller

import (
	"encoding/json"
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

/////////////////////////////////////////////////////////////
//                                                         //
//              DeletePermissionGroupController            //
//      This controller deletes an permission group        //
//    It requires deletePermissionGroupRequest instance    //
//                                                         //
/////////////////////////////////////////////////////////////
func DeletePermissionGroupController(c *fiber.Ctx) error {

	// init and parse the request object
	obj := deletePermissionGroupRequest{}
	err := json.Unmarshal(c.Body(), &obj)

	// check request
	if err != nil {
		if cfg, _ := config.ParseConfig(); cfg.ServerCFG.LogRequestErrors {
			utils.LogError(err.Error(), "DeletePermissionGroupController.go", 24)
		}

		res, _ := models.GetJSONResponse("Invaild JSON body", "#d41717", "error", "None", 200)
		return c.Send(res)
	}
	if !checkDeletePermissionGroupRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "#d41717", "ok", "None", 200)
		return c.Send(res)
	}

	// check login status
	if !actions.MysqlLoginWithToken(obj.Username, obj.Password, obj.Token) {
		res, _ := models.GetJSONResponse("You do not have the permission to perform this", "#d41717", "ok", "None", 200)
		return c.Send(res)
	}

	conn := actions.GetConn()
	defer conn.Close()

	// check if user has higher permission
	if actions.CheckUserHasHigherPermission(conn, obj.Username, 0, "permission."+obj.GroupName) {

		user := actions.GetUsersByPermission("%permission." + obj.GroupName + "%")

		for _, val := range user {

			split := strings.Split(val.Permissions, ";")

			// remove permission from user permission array
			reduced := utils.RemoveValueFromArray(split, "permission."+obj.GroupName)

			editedPerms := reduced[0]

			// build permission-string from reduced
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

		res, _ := models.GetJSONResponse("Successfully deleted PermissionGroup", "#1db004", "ok", "None", 200)
		return c.Send(res)
	}

	res, _ := models.GetJSONResponse("You do not have the permission to perform this", "#d41717", "ok", "None", 200)
	return c.Send(res)
}

// checks the request
// struct fields should not be default
func checkDeletePermissionGroupRequest(obj deletePermissionGroupRequest) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.GroupName != ""
}
