package permission_management

import (
	"github.com/MathisBurger/OpenInventory-Backend/controller/general"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	dbModels "github.com/MathisBurger/OpenInventory-Backend/database/models"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/gofiber/fiber/v2"
)

type listAllPermissionGroupsResponse struct {
	Message          string                     `json:"message"`
	PermissionGroups []dbModels.PermissionModel `json:"permission_groups"`
	Alert            string                     `json:"alert"`
}

////////////////////////////////////////////////////////////////////
//                                                                //
//                ListAllPermissionGroupsController               //
//           This controller fetches all permission groups        //
//        It requires models.LoginWithTokenRequest instance       //
//                                                                //
////////////////////////////////////////////////////////////////////
func ListAllPermissionGroupsController(c *fiber.Ctx) error {

	// init and parse the request object
	obj := models.LoginWithTokenRequest{
		Username: c.Query("username", ""),
		Password: c.Query("password", ""),
		Token:    c.Query("token", ""),
	}

	// check request
	if !general.CheckCheckCredsRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "#d41717", "ok", "None", 200)
		return c.Send(res)
	}

	// check login
	if !actions.MysqlLoginWithToken(obj.Username, obj.Password, obj.Token) {
		res, _ := models.GetJSONResponse("You do not have the permission to perform this", "#d41717", "ok", "None", 200)
		return c.Send(res)
	}

	return c.JSON(listAllPermissionGroupsResponse{
		"Successfully fetched all permission groups",
		actions.GetAllPermissions(),
		"#1db004",
	})
}
