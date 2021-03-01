package controller

import (
	"github.com/MathisBurger/OpenInventory-Backend/config"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

// ---------------------------------------------
//               deleteUserRequest
//    This struct contains login credentials
//                 and username
// ---------------------------------------------
type deleteUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
	User     string `json:"user"`
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
	obj := new(deleteUserRequest)
	err := c.BodyParser(obj)

	// check parsing error
	if err != nil {

		// log error if enabled
		if cfg, _ := config.ParseConfig(); cfg.ServerCFG.LogRequestErrors {
			utils.LogError(err.Error(), "DeleteUserController.go", 40)
		}

		// return "Invalid JSON body" response
		res, _ := models.GetJSONResponse("Invalid JSON body", "alert alert-danger", "error", "None", 200)
		return c.Send(res)
	}

	// check if request has been parsed correctly
	if !checkDeleteUserRequest(obj) {

		// send "Wrong JSON syntax" response
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}

	// check login status
	if actions.MysqlLoginWithTokenRoot(obj.Username, obj.Password, obj.Token) {

		// delete user
		actions.DeleteUser(obj.User)

		// return "Successfully deleted user" response
		resp, _ := models.GetJSONResponse("Successfully deleted user", "alert alert-success", "ok", "None", 200)
		return c.Send(resp)
	}

	// return "You do not have the permission to execute this" response
	res, _ := models.GetJSONResponse("You do not have the permission to execute this", "alert alert-danger", "ok", "None", 200)
	return c.Send(res)

}

/////////////////////////////////////////////////////////////
//                                                         //
//                 checkDeleteUserRequest                  //
//      This function is checking the request object       //
//        It requires the deleteUserRequest object         //
//                                                         //
/////////////////////////////////////////////////////////////
func checkDeleteUserRequest(obj *deleteUserRequest) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.User != ""
}
