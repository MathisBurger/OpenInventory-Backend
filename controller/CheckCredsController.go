package controller

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/config"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

/////////////////////////////////////////////////////////////
//                                                         //
//                   CheckCredsController                  //
//             This controller checks if username          //
//               password and token are valid              //
//                                                         //
/////////////////////////////////////////////////////////////
func CheckCredsController(c *fiber.Ctx) error {

	// init and parse the request object
	obj := models.LoginWithTokenRequest{}
	err := json.Unmarshal(c.Body(), &obj)

	// check request
	if err != nil {
		if cfg, _ := config.ParseConfig(); cfg.ServerCFG.LogRequestErrors {
			utils.LogError(err.Error(), "CheckCredsController.go", 17)
		}

		res, _ := models.GetJSONResponse("Invaild JSON body", "alert alert-danger", "error", "None", 200)
		return c.Send(res)
	}
	if !checkCheckCredsRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}

	// check login status
	if actions.MysqlLoginWithToken(obj.Username, obj.Password, obj.Token) {
		res, _ := models.GetJSONResponse("Login successful", "alert alert-success", "ok", "None", 200)
		return c.Send(res)
	}

	res, _ := models.GetJSONResponse("Login failed", "alert alert-warning", "ok", "None", 200)
	return c.Send(res)
}

// checks the request
// struct fields should not be default
func checkCheckCredsRequest(obj models.LoginWithTokenRequest) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != ""
}
