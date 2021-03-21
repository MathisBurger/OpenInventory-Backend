package permission_management

import (
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	dbModels "github.com/MathisBurger/OpenInventory-Backend/database/models"
	"github.com/MathisBurger/OpenInventory-Backend/models"
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
	obj := listAllPermGroupsOfTableRequest{
		Username:  c.Query("username", ""),
		Password:  c.Query("password", ""),
		Token:     c.Query("token", ""),
		TableName: c.Query("table_name", ""),
	}

	// check request
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
