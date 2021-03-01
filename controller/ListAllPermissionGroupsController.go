package controller

import (
	"github.com/MathisBurger/OpenInventory-Backend/config"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	dbModels "github.com/MathisBurger/OpenInventory-Backend/database/models"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

type listAllPermissionGroupsResponse struct {
	Message          string                     `json:"message"`
	PermissionGroups []dbModels.PermissionModel `json:"permission_groups"`
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
	obj := new(models.LoginWithTokenRequest)
	err := c.BodyParser(obj)

	// check request
	if err != nil {
		if cfg, _ := config.ParseConfig(); cfg.ServerCFG.LogRequestErrors {
			utils.LogError(err.Error(), "ListAllPermissionGroupsController.go", 21)
		}
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !checkCheckCredsRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}

	// check login
	if !actions.MysqlLoginWithToken(obj.Username, obj.Password, obj.Token) {
		res, _ := models.GetJSONResponse("You do not have the permission to perform this", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}

	return c.JSON(listAllPermissionGroupsResponse{
		"Successfully fetched all permission groups",
		actions.GetAllPermissions(),
	})
}
