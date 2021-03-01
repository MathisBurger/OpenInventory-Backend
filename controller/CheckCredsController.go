package controller

import (
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

	// initializing the request object
	obj := new(models.LoginWithTokenRequest)

	// parsing the body into the request object
	err := c.BodyParser(obj)

	// returns "Wrong JSON syntax" response if error is unequal nil
	if err != nil {

		// checks if request errors should be logged
		if cfg, _ := config.ParseConfig(); cfg.ServerCFG.LogRequestErrors {

			// log error
			utils.LogError(err.Error(), "CheckCredsController.go", 17)
		}

		// returns response
		res, _ := models.GetJSONResponse("Invaild JSON body", "alert alert-danger", "error", "None", 200)
		return c.Send(res)
	}

	// check if request has been parsed correctly
	if !checkCheckCredsRequest(obj) {

		// returns "Wrong JSON syntax" response
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}

	// check login status
	if actions.MysqlLoginWithToken(obj.Username, obj.Password, obj.Token) {

		// returns "Login successful" response if login successful
		res, _ := models.GetJSONResponse("Login successful", "alert alert-success", "ok", "None", 200)
		return c.Send(res)
	}

	// returns "Login failed" response if login failed
	res, _ := models.GetJSONResponse("Login failed", "alert alert-warning", "ok", "None", 200)
	return c.Send(res)
}

//////////////////////////////////////////////////////////
//                                                      //
//                checkCheckCredsRequest                //
//             consumes the request object              //
//   checks if struct fields are not the default value  //
//                                                      //
//////////////////////////////////////////////////////////
func checkCheckCredsRequest(obj *models.LoginWithTokenRequest) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != ""
}
