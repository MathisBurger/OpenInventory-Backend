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

type listAllPermsOfUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
	User     string `json:"user"`
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
	obj := listAllPermsOfUserRequest{}
	err := json.Unmarshal(c.Body(), &obj)

	// check request
	if err != nil {
		if cfg, _ := config.ParseConfig(); cfg.ServerCFG.LogRequestErrors {
			utils.LogError(err.Error(), "ListAllPermsOfUserController.go", 31)
		}
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "#d41717", "ok", "None", 200)
		return c.Send(res)
	}
	if !checkListAllPermsOfUserRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "#d41717", "ok", "None", 200)
		return c.Send(res)
	}

	// check login
	if !actions.MysqlLoginWithToken(obj.Username, obj.Password, obj.Token) {
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
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.User != ""
}