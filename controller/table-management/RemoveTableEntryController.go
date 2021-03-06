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

type removeTableEntryRequest struct {
	TableName string `json:"table_name"`
	RowID     int    `json:"row_id"`
}

////////////////////////////////////////////////////////////////////
//                                                                //
//                  RemoveTableEntryController                    //
//          This controller removes entry from table              //
//        It requires removeTableEntryRequest instance            //
//                                                                //
////////////////////////////////////////////////////////////////////
func RemoveTableEntryController(c *fiber.Ctx) error {

	// init and parse the request object
	obj := removeTableEntryRequest{}
	decrypted, err := e2e.DecryptBytes(c.Body())
	if err != nil {
		return c.SendStatus(400)
	}
	err = json.Unmarshal(decrypted, &obj)

	// check request
	if err != nil {
		if cfg, _ := config.ParseConfig(); cfg.ServerCFG.LogRequestErrors {
			utils.LogError(err.Error(), "RemoveTableEntryController.go", 25)
		}
		res, _ := models.GetJSONResponse("Invaild JSON body", "#d41717", "error", "None", 200)
		return c.Send(res)
	}
	if !checkRemoveTableEntryRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "#d41717", "ok", "None", 200)
		return c.Send(res)
	}

	// check login
	if ok, ident := middleware.ValidateAccessToken(c); ok {
		table := actions.GetTableByName(obj.TableName)
		conn := actions.GetConn()
		defer conn.Close()

		// check permission
		if actions.CheckUserHasHigherPermission(conn, ident, table.MinPermLvl, "") {

			// check deletion status
			if !actions.DeleteTableEntry(obj.RowID, obj.TableName) {
				res, _ := models.GetJSONResponse("EntryID does not exist", "alert alert-warning", "ok", "None", 200)
				return c.Send(res)
			}

			actions.ChangeNumOfEntrysBy(obj.TableName, -1)

			res, _ := models.GetJSONResponse("Successfully deleted entry", "#1db004", "ok", "None", 200)
			return c.Send(res)
		}
		res, _ := models.GetJSONResponse("You do not have the permission perform this", "#d41717", "ok", "None", 200)
		return c.Send(res)
	}
	res, _ := models.GetJSONResponse("You do not have the permission perform this", "#d41717", "ok", "None", 200)
	return c.Send(res)
}

// checks the request
// struct fields should not be default
func checkRemoveTableEntryRequest(obj removeTableEntryRequest) bool {
	return obj.TableName != "" && obj.RowID > 0
}
