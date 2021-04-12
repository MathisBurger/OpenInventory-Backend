package table_management

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/e2e"

	"github.com/MathisBurger/OpenInventory-Backend/config"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/middleware"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

type deleteTableRequest struct {
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

	// init and parse the request object
	obj := deleteTableRequest{}
	decrypted, err := e2e.DecryptBytes(c.Body())
	if err != nil {
		return c.SendStatus(400)
	}
	err = json.Unmarshal(decrypted, &obj)

	// check request
	if err != nil {
		if cfg, _ := config.ParseConfig(); cfg.ServerCFG.LogRequestErrors {
			utils.LogError(err.Error(), "DeleteTableController.go", 23)
		}

		res, _ := models.GetJSONResponse("Wrong JSON syntax", "#d41717", "ok", "None", 200)
		return c.Send(res)
	}
	if !checkDeleteTableRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "#d41717", "ok", "None", 200)
		return c.Send(res)
	}

	// check login status
	if ok, ident := middleware.ValidateAccessToken(c); ok {

		conn := actions.GetConn()
		defer conn.Close()

		table := actions.GetTableByName(obj.TableName)

		// check permission of user
		if actions.CheckUserHasHigherPermission(conn, ident, table.MinPermLvl, "") {

			actions.DropTable(obj.TableName)

			res, _ := models.GetJSONResponse("Successfully deleted table", "#1db004", "ok", "None", 200)
			return c.Send(res)
		}

		res, _ := models.GetJSONResponse("You do not have the permission to perform this", "#d41717", "ok", "None", 200)
		return c.Send(res)

	}

	res, _ := models.GetJSONResponse("You do not have the permission to perform this", "#d41717", "ok", "None", 200)
	return c.Send(res)

}

// checks the request
// struct fields should not be default
func checkDeleteTableRequest(obj deleteTableRequest) bool {
	return obj.TableName != ""
}
