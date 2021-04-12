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

type addTableEntryRequest struct {
	TableName string                 `json:"table_name"`
	Row       map[string]interface{} `json:"row"`
}

////////////////////////////////////////////////////////////////
//                                                            //
//                   AddTableEntryController                  //
//   This controller adds an table entry to the given table   //
//       It requires login credentials and the table-name     //
//                                                            //
////////////////////////////////////////////////////////////////
func AddTableEntryController(c *fiber.Ctx) error {

	// init and parse the request object
	obj := addTableEntryRequest{}
	decrypted, err := e2e.DecryptBytes(c.Body())
	if err != nil {
		return c.SendStatus(400)
	}
	err = json.Unmarshal(decrypted, &obj)

	// check request
	if err != nil {

		if cfg, _ := config.ParseConfig(); cfg.ServerCFG.LogRequestErrors {
			utils.LogError(err.Error(), "AddTableEntryController.go", 22)
		}

		res, _ := models.GetJSONResponse("Wrong JSON syntax", "#d41717", "ok", "None", 200)
		return c.Send(res)
	}
	if !checkAddTableEntryRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "#d41717", "ok", "None", 200)
		return c.Send(res)
	}

	// check adding status
	if ok, ident := middleware.ValidateAccessToken(c); ok {
		actions.AddTableEntry(ident, obj.TableName, obj.Row)
		res, _ := models.GetJSONResponse("successful", "#1db004", "ok", "None", 200)
		return c.Send(res)
	}

	res, _ := models.GetJSONResponse("creation failed", "#d41717", "ok", "None", 200)
	return c.Send(res)
}

// checks the request
// struct fields should not be default
func checkAddTableEntryRequest(obj addTableEntryRequest) bool {
	return obj.TableName != "" && len(obj.Row) > 0
}
