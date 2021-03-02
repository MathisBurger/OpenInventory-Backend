package general

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/config"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

////////////////////////////////////////////////////////////////////
//                                                                //
//                        LoginController                         //
//               This controller executes user login              //
//             It requires models.LoginRequest instance           //
//                                                                //
////////////////////////////////////////////////////////////////////
func LoginController(c *fiber.Ctx) error {

	// init and parse the request object
	obj := models.LoginRequest{}
	err := json.Unmarshal(c.Body(), &obj)

	// check request
	if err != nil {
		if cfg, _ := config.ParseConfig(); cfg.ServerCFG.LogRequestErrors {
			utils.LogError(err.Error(), "LoginController.go", 16)
		}
		res, _ := models.GetJSONResponse("Invaild JSON body", "alert alert-danger", "error", "None", 200)
		return c.Send(res)
	}
	if !checkLoginRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}

	// check login
	status, token := actions.MysqlLogin(obj.Username, obj.Password)
	if status {
		res, _ := models.GetJSONResponse("Login successful", "#1db004", "ok", token, 200)
		return c.Send(res)
	}
	res, _ := models.GetJSONResponse("Login failed", "#d41717", "ok", "None", 200)
	return c.Send(res)

}

func checkLoginRequest(obj models.LoginRequest) bool {
	return obj.Username != "" && obj.Password != ""
}
