package permission_management

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/e2e"
	"strings"

	"github.com/MathisBurger/OpenInventory-Backend/config"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/middleware"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

type deletePermissionGroupRequest struct {
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
	decrypted, err := e2e.DecryptBytes(c.Body())
	if err != nil {
		return c.SendStatus(400)
	}
	err = json.Unmarshal(decrypted, &obj)

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
	if ok, ident := middleware.ValidateAccessToken(c); ok {

		conn := actions.GetConn()
		defer conn.Close()

		// check if user has higher permission
		if actions.CheckUserHasHigherPermission(conn, ident, 0, "permission."+obj.GroupName) {

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

	res, _ := models.GetJSONResponse("You do not have the permission to perform this", "#d41717", "ok", "None", 200)
	return c.Send(res)
}

// checks the request
// struct fields should not be default
func checkDeletePermissionGroupRequest(obj deletePermissionGroupRequest) bool {
	return obj.GroupName != ""
}
