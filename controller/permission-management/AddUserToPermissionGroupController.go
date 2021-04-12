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

type addUserToPermissionGroupRequest struct {
	Permission string `json:"permission"`
	User       string `json:"user"`
}

/////////////////////////////////////////////////////////////
//                                                         //
//            AddUserToPermissionGroupController           //
//          This controller adds permission to user        //
//       It requires login credentials, permission-name    //
//                      and username                       //
//                                                         //
/////////////////////////////////////////////////////////////
func AddUserToPermissionGroupController(c *fiber.Ctx) error {

	// init and parse the request object
	obj := addUserToPermissionGroupRequest{}
	decrypted, err := e2e.DecryptBytes(c.Body())
	if err != nil {
		return c.SendStatus(400)
	}
	err = json.Unmarshal(decrypted, &obj)

	// check request
	if err != nil {
		if cfg, _ := config.ParseConfig(); cfg.ServerCFG.LogRequestErrors {
			utils.LogError(err.Error(), "AddUserToPermissionGroupController.go", 26)
		}

		res, _ := models.GetJSONResponse("Wrong JSON syntax", "#d41717", "ok", "None", 200)
		return c.Send(res)
	}
	if !checkAddUserToPermissionGroupRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "#d41717", "ok", "None", 200)
		return c.Send(res)
	}

	// check login status
	if ok, intent := middleware.ValidateAccessToken(c); ok {
		conn := actions.GetConn()
		defer conn.Close()

		// check permission
		if actions.CheckUserHasHigherPermission(conn, intent, 0, obj.Permission) {
			userexists, user := actions.GetUserByUsername(obj.User)

			if !userexists {
				res, _ := models.GetJSONResponse("User does not exist", "#d41717", "ok", "None", 200)
				return c.Send(res)
			}

			// check if user already owns this permission
			if utils.ContainsStr(strings.Split(user.Permissions, ";"), obj.Permission) {

				res, _ := models.GetJSONResponse("The user is already member of this group", "#d41717", "ok", "None", 200)
				return c.Send(res)
			}

			// check if permission exists
			if permexists, _ := actions.GetPermissionByName(obj.Permission); !permexists {

				res, _ := models.GetJSONResponse("This permissiongroup does not exist", "#d41717", "ok", "None", 200)
				return c.Send(res)
			}

			finalPermissions := user.Permissions + ";" + obj.Permission

			actions.UpdateUserPermission(obj.User, finalPermissions)

			if obj.Permission == "default.root" {
				actions.UpdateUserRoor(true, obj.User)
			}

			res, _ := models.GetJSONResponse("User added to permissiongroup", "#1db004", "ok", "None", 200)
			return c.Send(res)
		}

		res, _ := models.GetJSONResponse("Your permission-level is too low", "#d41717", "ok", "None", 200)
		return c.Send(res)
	}

	res, _ := models.GetJSONResponse("Wrong login credentials", "#d41717", "ok", "None", 200)
	return c.Send(res)
}

// checks the request
// struct fields should not be default
func checkAddUserToPermissionGroupRequest(obj addUserToPermissionGroupRequest) bool {
	return obj.Permission != "" && obj.User != ""
}
