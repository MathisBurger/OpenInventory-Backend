package controller

import (
	"github.com/MathisBurger/OpenInventory-Backend/config"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

// ---------------------------------------------
//              deleteTableRequest
//    This struct contains login credentials
//                and table-name
// ---------------------------------------------
type deleteTableRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Token     string `json:"token"`
	TableName string `json:"table_name"`
}

/////////////////////////////////////////////////////////////
//                                                         //
//                 DeleteTableController                   //
//            This controller deletes an table             //
//         It requires deleteTableRequest instance         //
//                                                         //
/////////////////////////////////////////////////////////////
func DeleteTableController(c *fiber.Ctx) error {

	// initializing the request object
	obj := new(deleteTableRequest)

	// parsing the body into the request object
	err := c.BodyParser(obj)

	// returns "Wrong JSON syntax" response if error is unequal nil
	if err != nil {

		// checks if request errors should be logged
		if cfg, _ := config.ParseConfig(); cfg.ServerCFG.LogRequestErrors {

			// log error
			utils.LogError(err.Error(), "DeleteTableController.go", 23)
		}

		// returns response
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}

	// check if request has been parsed correctly
	if !checkDeleteTableRequest(obj) {

		// returns "Wrong JSON syntax" response
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}

	// check login status
	if actions.MysqlLoginWithToken(obj.Username, obj.Password, obj.Token) {

		// get connection
		conn := actions.GetConn()
		defer conn.Close()

		table := actions.GetTableByName(obj.TableName)

		// check if user has higher permission
		if actions.CheckUserHasHigherPermission(conn, obj.Username, table.MinPermLvl, "") {

			// delete table
			actions.DropTable(obj.TableName)

			// send response
			res, _ := models.GetJSONResponse("Successfully deleted table", "alert alert-success", "ok", "None", 200)
			return c.Send(res)
		}

		// send invalid permission response
		res, _ := models.GetJSONResponse("You do not have the permission to perform this", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)

	}

	// send invalid permission response
	res, _ := models.GetJSONResponse("You do not have the permission to perform this", "alert alert-danger", "ok", "None", 200)
	return c.Send(res)

}

/////////////////////////////////////////////////////////////
//                                                         //
//                 checkDeleteTableRequest                 //
//      This function is checking the request object       //
//        It requires the deleteTableRequest object        //
//                                                         //
/////////////////////////////////////////////////////////////
func checkDeleteTableRequest(obj *deleteTableRequest) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.TableName != ""
}
