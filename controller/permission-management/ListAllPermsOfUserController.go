package permission_management

import (
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	dbModels "github.com/MathisBurger/OpenInventory-Backend/database/models"
	"github.com/MathisBurger/OpenInventory-Backend/middleware"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/gofiber/fiber/v2"
)

type listAllPermsOfUserRequest struct {
	User string `json:"user"`
}

type listAllPermsOfUserResponse struct {
	Message     string                     `json:"message"`
	Permissions []dbModels.PermissionModel `json:"permissions"`
	Status      string                     `json:"status"`
	HttpStatus  int                        `json:"http_status"`
	Alert       string                     `json:"alert"`
}

////////////////////////////////////////////////////////////////////
//                                                                //
//                  ListAllPermOfUserController                   //
//    This controller fetches all permissions of user             //
//         It requires listAllPermsOfUserRequest instance         //
//                                                                //
////////////////////////////////////////////////////////////////////
func ListAllPermOfUserController(c *fiber.Ctx) error {

	// init and parse the request object
	obj := listAllPermsOfUserRequest{
		User: c.Query("user", ""),
	}

	// check request
	if !checkListAllPermsOfUserRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "#d41717", "ok", "None", 200)
		return c.Send(res)
	}

	// check login
	if ok, _ := middleware.ValidateAccessToken(c); !ok {
		res, _ := models.GetJSONResponse("You do not have the permission to perform this", "#d41717", "ok", "None", 200)
		return c.Send(res)
	}

	return c.JSON(listAllPermsOfUserResponse{
		"Successfully fetched all user permissions",
		actions.GetPermissionsOfUser(obj.User),
		"ok",
		200,
		"#1db004",
	})
}

// checks the request
// struct fields should not be default
func checkListAllPermsOfUserRequest(obj listAllPermsOfUserRequest) bool {
	return obj.User != ""
}
