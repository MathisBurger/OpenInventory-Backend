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
//         createPermissionGroupRequest
//    This struct contains login credentials
//         and permission information
// ---------------------------------------------
type createPermissionGroupRequest struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	Token          string `json:"token"`
	PermissionInfo struct {
		Name            string `json:"name"`
		ColorCode       string `json:"color_code"`
		PermissionLevel int    `json:"permission_level"`
	} `json:"permission_info"`
}

/////////////////////////////////////////////////////////////
//                                                         //
//               CreatePermissionGroupController           //
//          This controller creates a new permission       //
//     It requires login credentials and permission-info   //
//                                                         //
/////////////////////////////////////////////////////////////
func CreatePermissionGroupController(c *fiber.Ctx) error {

	// initializing the request object
	obj := new(createPermissionGroupRequest)

	// parsing the body into the request object
	err := c.BodyParser(obj)

	// returns "Wrong JSON syntax" response if error is unequal nil
	if err != nil {

		// checks if request errors should be logged
		if cfg, _ := config.ParseConfig(); cfg.ServerCFG.LogRequestErrors {

			// log error
			utils.LogError(err.Error(), "CreatePermissionGroupController.go", 29)
		}

		// returns response
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}

	// check if request has been parsed correctly
	if !checkCreatePermissionGroupRequest(obj) {

		// returns "Wrong JSON syntax" response
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}

	// check login status
	if !actions.MysqlLoginWithToken(obj.Username, obj.Password, obj.Token) {

		// returns "You do not have the permission to perform this" response if login failed
		res, _ := models.GetJSONResponse("You do not have the permission to perform this", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}

	// get permission group input
	permGroupInputStatus := checkPermissionGroupInput(obj)

	// check if permission permGroupInputStatus is nil
	if permGroupInputStatus != nil {

		// send []byte permGroupInputStatus as response
		return c.Send(permGroupInputStatus)
	}

	// check if permission exists
	if exists, _ := actions.GetPermissionByName(obj.PermissionInfo.Name); exists {

		// return "This group already exists" response if group already exists
		res, _ := models.GetJSONResponse("This group already exists", "alert alert-warning", "ok", "None", 200)
		return c.Send(res)
	}

	// create permission group
	actions.InsertPermissionGroup(obj.PermissionInfo.Name, obj.PermissionInfo.ColorCode, obj.PermissionInfo.PermissionLevel)

	// return "Created permission-group" response
	res, _ := models.GetJSONResponse("Created permission-group", "alert alert-success", "ok", "None", 200)
	return c.Send(res)
}

//////////////////////////////////////////////////////////
//                                                      //
//          checkCreatePermissionGroupRequest           //
//             consumes the request object              //
//   checks if struct fields are not the default value  //
//                                                      //
//////////////////////////////////////////////////////////
func checkCreatePermissionGroupRequest(obj *createPermissionGroupRequest) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.PermissionInfo.Name != "" && obj.PermissionInfo.ColorCode != "" && obj.PermissionInfo.PermissionLevel > 0
}

//////////////////////////////////////////////////////////
//                                                      //
//              checkPermissionGroupInput               //
//             consumes the request object              //
//            checks if the permission input            //
//                                                      //
//////////////////////////////////////////////////////////
func checkPermissionGroupInput(obj *createPermissionGroupRequest) []byte {

	// check if permission-name contains ';'
	if strings.Contains(obj.PermissionInfo.Name, ";") {

		// returns "';' is not allowed in group name" response if permission-name contains ';'
		res, _ := models.GetJSONResponse("';' is not allowed in group name", "alert alert-danger", "ok", "None", 200)
		return res
	}

	// create connection
	conn := actions.GetConn()

	// defer close connection
	defer conn.Close()

	// check if user has higher permission
	if !actions.CheckUserHasHigherPermission(conn, obj.Username, obj.PermissionInfo.PermissionLevel, "") {

		// returns "Your permission is not high enough" if permission if not high enough
		res, _ := models.GetJSONResponse("Your permission is not high enough", "alert alert-danger", "ok", "None", 200)
		return res
	}

	// returns nil if everything is okay
	return nil
}
