package user_management

import (
	"encoding/json"

	"github.com/MathisBurger/OpenInventory-Backend/config"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/middleware"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

type deleteUserRequest struct {
	User string `json:"user"`
}

/////////////////////////////////////////////////////////////
//                                                         //
//                 DeleteUserController                    //
//    This controller deletes the user given in request    //
//        It requires deleteUserRequest instance           //
//                                                         //
/////////////////////////////////////////////////////////////
func DeleteUserController(c *fiber.Ctx) error {

	// init and parse the request object
	obj := deleteUserRequest{}
	err := json.Unmarshal(c.Body(), &obj)

	// check request
	if err != nil {
		if cfg, _ := config.ParseConfig(); cfg.ServerCFG.LogRequestErrors {
			utils.LogError(err.Error(), "DeleteUserController.go", 40)
		}

		res, _ := models.GetJSONResponse("Invalid JSON body", "#d41717", "error", "None", 200)
		return c.Send(res)
	}
	if !checkDeleteUserRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "#d41717", "ok", "None", 200)
		return c.Send(res)
	}

	conn := actions.GetConn()
	defer conn.Close()

	// check login status
	if ok, ident := middleware.ValidateAccessToken(c); ok &&
		actions.CheckUserHasHigherPermission(conn, ident, actions.GetHighestPermission(conn, obj.User), "") {

		actions.DeleteUser(obj.User)

		resp, _ := models.GetJSONResponse("Successfully deleted user", "#1db004", "ok", "None", 200)
		return c.Send(resp)
	}

	res, _ := models.GetJSONResponse("You do not have the permission to execute this", "#d41717", "ok", "None", 200)
	return c.Send(res)

}

// checks the request
// struct fields should not be default
func checkDeleteUserRequest(obj deleteUserRequest) bool {
	return obj.User != ""
}
