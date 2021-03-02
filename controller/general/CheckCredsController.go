package general

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

		res, _ := models.GetJSONResponse("Invaild JSON body", "#d41717", "error", "None", 200)
		return c.Send(res)
	}
	if !CheckCheckCredsRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "#d41717", "ok", "None", 200)
		return c.Send(res)
	}

	// check login status
	if actions.MysqlLoginWithToken(obj.Username, obj.Password, obj.Token) {
		res, _ := models.GetJSONResponse("Login successful", "#1db004", "ok", "None", 200)
		return c.Send(res)
	}

	res, _ := models.GetJSONResponse("Login failed", "#d41717", "ok", "None", 200)
	return c.Send(res)
}

// checks the request
// struct fields should not be default
func CheckCheckCredsRequest(obj models.LoginWithTokenRequest) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != ""
}
