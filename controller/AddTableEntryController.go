package controller

import (
	"github.com/MathisBurger/OpenInventory-Backend/config"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

// ---------------------------------------------
//         addTableEntryRequest
//  This struct contains all request parameter
//     It is parsed in controller below
// ---------------------------------------------
type addTableEntryRequest struct {
	Username  string                 `json:"username"`
	Password  string                 `json:"password"`
	Token     string                 `json:"token"`
	TableName string                 `json:"table_name"`
	Row       map[string]interface{} `json:"row"`
}

////////////////////////////////////////////////////////////////
//                                                            //
//                   AddTableEntryController                  //
//   This controller adds an table entry to the given table   //
//       It requires login credentials and the tablename      //
//                                                            //
////////////////////////////////////////////////////////////////
func AddTableEntryController(c *fiber.Ctx) error {

	// initializing the request object
	obj := new(addTableEntryRequest)

	// parsing the body into the request object
	err := c.BodyParser(obj)

	// returns "Wrong JSON syntax" response if error is unequal nil
	if err != nil {

		// checks if request errors should be logged
		if cfg, _ := config.ParseConfig(); cfg.ServerCFG.LogRequestErrors {

			// log error
			utils.LogError(err.Error(), "AddTableEntryController.go", 22)

		}

		// returns response
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}

	// check if request has been parsed correctly
	if !checkAddTableEntryRequest(obj) {

		// returns "Wrong JSON syntax" response
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)

	}

	// try to insert the table entry
	status := actions.AddTableEntry(obj.Username, obj.Password, obj.Token, obj.TableName, obj.Row)

	// checks status of table insertion
	if status {

		// returns "successful" response if status is true
		res, _ := models.GetJSONResponse("successful", "alert alert-success", "ok", "None", 200)
		return c.Send(res)
	}

	// returns "creation failed" response if status is false
	res, _ := models.GetJSONResponse("creation failed", "alert alert-danger", "ok", "None", 200)
	return c.Send(res)
}

/////////////////////////////////////////////////////////
//                                                     //
//             checkAddTableEntryRequest               //
//  consumes the request object of the upper function  //
//  checks if struct fields are not the default value  //
//                                                     //
/////////////////////////////////////////////////////////
func checkAddTableEntryRequest(obj *addTableEntryRequest) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.TableName != "" && len(obj.Row) > 0
}
