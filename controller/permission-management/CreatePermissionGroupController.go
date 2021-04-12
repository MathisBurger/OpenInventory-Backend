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

type createPermissionGroupRequest struct {
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

	// init and parse the request object
	obj := createPermissionGroupRequest{}
	decrypted, err := e2e.DecryptBytes(c.Body())
	if err != nil {
		return c.SendStatus(400)
	}
	err = json.Unmarshal(decrypted, &obj)

	// check request
	if err != nil {
		if cfg, _ := config.ParseConfig(); cfg.ServerCFG.LogRequestErrors {
			utils.LogError(err.Error(), "CreatePermissionGroupController.go", 29)
		}

		res, _ := models.GetJSONResponse("Wrong JSON syntax", "#d41717", "ok", "None", 200)
		return c.Send(res)
	}
	if !checkCreatePermissionGroupRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "#d41717", "ok", "None", 200)
		return c.Send(res)
	}

	// check login status
	if ok, ident := middleware.ValidateAccessToken(c); ok {

		permGroupInputStatus := checkPermissionGroupInput(obj, ident)

		// check if request contains invalid parameter
		if permGroupInputStatus != nil {
			return c.Send(permGroupInputStatus)
		}

		if exists, _ := actions.GetPermissionByName(obj.PermissionInfo.Name); exists {
			res, _ := models.GetJSONResponse("This group already exists", "#d41717", "ok", "None", 200)
			return c.Send(res)
		}

		actions.InsertPermissionGroup(obj.PermissionInfo.Name, obj.PermissionInfo.ColorCode, obj.PermissionInfo.PermissionLevel)

		res, _ := models.GetJSONResponse("Created permission-group", "#1db004", "ok", "None", 200)
		return c.Send(res)
	}

	res, _ := models.GetJSONResponse("You do not have the permission to perform this", "#d41717", "ok", "None", 200)
	return c.Send(res)
}

// checks the request
// struct fields should not be default
func checkCreatePermissionGroupRequest(obj createPermissionGroupRequest) bool {
	return obj.PermissionInfo.Name != "" && obj.PermissionInfo.ColorCode != "" && obj.PermissionInfo.PermissionLevel > 0
}

// checks for disallowed syntax in createPermissionGroupRequest object
// returns a []byte response which can be send as response
func checkPermissionGroupInput(obj createPermissionGroupRequest, username string) []byte {

	// check if permission-name contains ';'
	if strings.Contains(obj.PermissionInfo.Name, ";") {

		// returns "';' is not allowed in group name" response if permission-name contains ';'
		res, _ := models.GetJSONResponse("';' is not allowed in group name", "#d41717", "ok", "None", 200)
		return res
	}

	conn := actions.GetConn()
	defer conn.Close()

	// check if user has higher permission
	if !actions.CheckUserHasHigherPermission(conn, username, obj.PermissionInfo.PermissionLevel, "") {

		res, _ := models.GetJSONResponse("Your permission is not high enough", "#d41717", "ok", "None", 200)
		return res
	}

	return nil
}
