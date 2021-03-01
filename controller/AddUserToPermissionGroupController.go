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
//        addUserToPermissionGroupRequest
//    This struct contains login credentials,
//         permission name and username
// ---------------------------------------------
type addUserToPermissionGroupRequest struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	Token      string `json:"token"`
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

	// initializing the request object
	obj := new(addUserToPermissionGroupRequest)

	// parsing the body into the request object
	err := c.BodyParser(obj)

	// returns "Wrong JSON syntax" response if error is unequal nil
	if err != nil {

		// checks if request errors should be logged
		if cfg, _ := config.ParseConfig(); cfg.ServerCFG.LogRequestErrors {

			// log error
			utils.LogError(err.Error(), "AddUserToPermissionGroupController.go", 26)
		}

		// returns response
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}

	// check if request has been parsed correctly
	if !checkAddUserToPermissionGroupRequest(obj) {

		// returns "Wrong JSON syntax" response
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}

	// check login status
	if !actions.MysqlLoginWithToken(obj.Username, obj.Password, obj.Token) {

		// returns "Wrong login credentials" response if login failed
		res, _ := models.GetJSONResponse("Wrong login credentials", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}

	// get connection
	conn := actions.GetConn()

	// defer close connection
	defer conn.Close()

	// check if logged in user has higher permission then user (new permission)
	if actions.CheckUserHasHigherPermission(conn, obj.Username, 0, obj.Permission) {

		// get user existence and instance of "models.UserModel" struct by username
		userexists, user := actions.GetUserByUsername(obj.User)

		// if user does not exist
		if !userexists {

			// returns "User does not exist" response
			res, _ := models.GetJSONResponse("User does not exist", "alert alert-danger", "ok", "None", 200)
			return c.Send(res)
		}

		// check if user already owns this permission
		if utils.ContainsStr(strings.Split(user.Permissions, ";"), obj.Permission) {

			// returns "The user is already member of this group" response
			res, _ := models.GetJSONResponse("The user is already member of this group", "alert alert-warning", "ok", "None", 200)
			return c.Send(res)
		}

		// check if permission exists
		if permexists, _ := actions.GetPermissionByName(obj.Permission); !permexists {

			// returns "This permissiongroup does not exist" response if permission does not exists
			res, _ := models.GetJSONResponse("This permissiongroup does not exist", "alert alert-warning", "ok", "None", 200)
			return c.Send(res)
		}

		// get final permission strings
		finalPermissions := user.Permissions + ";" + obj.Permission

		// update user permission
		actions.UpdateUserPermission(obj.User, finalPermissions)

		// returns "User added to permissiongroup" response
		res, _ := models.GetJSONResponse("User added to permissiongroup", "alert alert-success", "ok", "None", 200)
		return c.Send(res)
	}

	// returns "Your permission-level is too low" response if permission is too low
	res, _ := models.GetJSONResponse("Your permission-level is too low", "alert alert-warning", "ok", "None", 200)
	return c.Send(res)
}

//////////////////////////////////////////////////////////
//                                                      //
//          checkAddUserToPermissionGroupRequest        //
//             consumes the request object              //
//   checks if struct fields are not the default value  //
//                                                      //
//////////////////////////////////////////////////////////
func checkAddUserToPermissionGroupRequest(obj *addUserToPermissionGroupRequest) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.Permission != "" && obj.User != ""
}
