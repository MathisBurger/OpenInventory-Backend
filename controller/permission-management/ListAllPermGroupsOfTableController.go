package permission_management

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/config"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	dbModels "github.com/MathisBurger/OpenInventory-Backend/database/models"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

type listAllPermGroupsOfTableRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Token     string `json:"token"`
	TableName string `json:"table_name"`
}

type listAllPermGroupsOfTableResponse struct {
	PermGroups []dbModels.PermissionModel `json:"perm_groups"`
	Message    string                     `json:"message"`
	Alert      string                     `json:"alert"`
}

////////////////////////////////////////////////////////////////////
//                                                                //
//               ListAllPermGroupsOfTableController               //
//      This controller fetches all permission groups of table    //
//       It requires listAllPermGroupsOfTableRequest instance     //
//                                                                //
////////////////////////////////////////////////////////////////////
func ListAllPermGroupsOfTableController(c *fiber.Ctx) error {

	// init and parse the request object
	obj := listAllPermGroupsOfTableRequest{}
	err := json.Unmarshal(c.Body(), &obj)

	// check request
	if err != nil {
		if cfg, _ := config.ParseConfig(); cfg.ServerCFG.LogRequestErrors {
			utils.LogError(err.Error(), "ListAllPermGroupsOfTableController.go", 29)
		}
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "#d41717", "ok", "None", 200)
		return c.Send(res)
	}
	if !checkListAllPermGroupsOfTableRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "#d41717", "ok", "None", 200)
		return c.Send(res)
	}

	// check login
	if !actions.MysqlLoginWithToken(obj.Username, obj.Password, obj.Token) {
		res, _ := models.GetJSONResponse("You do not have the permission to perform this", "#d41717", "ok", "None", 200)
		return c.Send(res)
	}

	conn := actions.GetConn()
	defer conn.Close()

	table := actions.GetTableByName(obj.TableName)

	// check permission
	if !actions.CheckUserHasHigherPermission(conn, obj.Username, table.MinPermLvl, "") {
		res, _ := models.GetJSONResponse("Your permission is not high enough to view this table", "#d41717", "ok", "None", 200)
		return c.Send(res)
	}

	return c.JSON(listAllPermGroupsOfTableResponse{
		actions.GetAllPermissionsWithHigherPermLvl(table.MinPermLvl),
		"Successfully fetched all permissiongroups of table",
		"#1db004",
	})
}

// checks the request
// struct fields should not be default
func checkListAllPermGroupsOfTableRequest(obj listAllPermGroupsOfTableRequest) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.TableName != ""
}
