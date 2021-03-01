package controller

import (
	"github.com/MathisBurger/OpenInventory-Backend/config"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
	"strings"
)

// ---------------------------------------------
//         deletePermissionGroupRequest
//    This struct contains login credentials
//                and group-name
// ---------------------------------------------
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

	// initializing the request object
	obj := new(deletePermissionGroupRequest)

	// parsing the body into the request object
	err := c.BodyParser(obj)

	// returns "Wrong JSON syntax" response if error is unequal nil
	if err != nil {

		// checks if request errors should be logged
		if cfg, _ := config.ParseConfig(); cfg.ServerCFG.LogRequestErrors {

			// log error
			utils.LogError(err.Error(), "DeletePermissionGroupController.go", 24)
		}

		// returns response
		res, _ := models.GetJSONResponse("Invaild JSON body", "alert alert-danger", "error", "None", 200)
		return c.Send(res)
	}

	// check if request has been parsed correctly
	if !checkDeletePermissionGroupRequest(obj) {

		// returns "Wrong JSON syntax" response
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}

	// check login status
	if !actions.MysqlLoginWithToken(obj.Username, obj.Password, obj.Token) {

		// returns "You do not have the permission to perform this" response
		res, _ := models.GetJSONResponse("You do not have the permission to perform this", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}

	// get connection
	conn := actions.GetConn()

	// defer close connection
	defer conn.Close()

	// check if user has higher permission
	if actions.CheckUserHasHigherPermission(conn, obj.Username, 0, "permission."+obj.GroupName) {

		// get all user with given permission
		user := actions.GetUsersByPermission("%permission." + obj.GroupName + "%")

		// iterate trough user array
		for _, val := range user {

			// split permission-string to single permissions
			split := strings.Split(val.Permissions, ";")

			// remove permission from user permission array
			reduced := utils.RemoveValueFromArray(split, "permission."+obj.GroupName)

			// get edited permissions from reduced permission array
			editedPerms := reduced[0]

			// build permission-string from reduced
			for i, val2 := range reduced {
				if i == 0 {
					continue
				}
				editedPerms += ";" + val2
			}

			// update user permission
			actions.UpdateUserPermission(val.Username, editedPerms)
		}

		// create sql statement for deleting permission from table
		stmt, _ := conn.Prepare("DELETE FROM `inv_permissions` WHERE `name`=?")

		// defer close statement
		defer stmt.Close()

		// execute statement and fetch error
		_, err = stmt.Exec("permission." + obj.GroupName)

		// check if error != nil
		if err != nil {

			// log error
			utils.LogError(err.Error(), "DeletePermissionGroupController.go", 55)
		}

		// return "Successfully deleted PermissionGroup" response
		res, _ := models.GetJSONResponse("Successfully deleted PermissionGroup", "alert alert-success", "ok", "None", 200)
		return c.Send(res)
	}

	// return "You do not have the permission to perform this" response
	res, _ := models.GetJSONResponse("You do not have the permission to perform this", "alert alert-danger", "ok", "None", 200)
	return c.Send(res)
}

/////////////////////////////////////////////////////////////
//                                                         //
//             checkDeletePermissionGroupRequest           //
//      This function is checking the request object       //
//   It requires the deletePermissionGroupRequest object   //
//                                                         //
/////////////////////////////////////////////////////////////
func checkDeletePermissionGroupRequest(obj *deletePermissionGroupRequest) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.GroupName != ""
}
