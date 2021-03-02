package controller

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/config"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

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
//       It requires login credentials and the table-name     //
//                                                            //
////////////////////////////////////////////////////////////////
func AddTableEntryController(c *fiber.Ctx) error {

	// init and parse the request object
	obj := addTableEntryRequest{}
	err := json.Unmarshal(c.Body(), &obj)

	// check request
	if err != nil {

		if cfg, _ := config.ParseConfig(); cfg.ServerCFG.LogRequestErrors {
			utils.LogError(err.Error(), "AddTableEntryController.go", 22)
		}

		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !checkAddTableEntryRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}

	// check adding status
	if actions.AddTableEntry(obj.Username, obj.Password, obj.Token, obj.TableName, obj.Row) {
		res, _ := models.GetJSONResponse("successful", "alert alert-success", "ok", "None", 200)
		return c.Send(res)
	}

	res, _ := models.GetJSONResponse("creation failed", "alert alert-danger", "ok", "None", 200)
	return c.Send(res)
}

// checks the request
// struct fields should not be default
func checkAddTableEntryRequest(obj addTableEntryRequest) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.TableName != "" && len(obj.Row) > 0
}
