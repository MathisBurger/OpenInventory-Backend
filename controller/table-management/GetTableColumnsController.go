package table_management

import (
	"encoding/json"
	"fmt"
	"github.com/MathisBurger/OpenInventory-Backend/config"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

type getTableColumnsRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Token     string `json:"token"`
	TableName string `json:"table_name"`
}

type getTableColumnsResponse struct {
	Message string      `json:"message"`
	Alert   string      `json:"alert"`
	Columns interface{} `json:"columns"`
}

////////////////////////////////////////////////////////////////////
//                                                                //
//                   GetTableColumnsController                    //
//          This controller fetches all columns of table          //
//          It requires getTableColumnsRequest instance           //
//                                                                //
////////////////////////////////////////////////////////////////////
func GetTableColumnsController(c *fiber.Ctx) error {

	// init and parse the request object
	obj := getTableColumnsRequest{}
	err := json.Unmarshal(c.Body(), &obj)

	// check request
	if err != nil {
		if cfg, _ := config.ParseConfig(); cfg.ServerCFG.LogRequestErrors {
			utils.LogError(err.Error(), "GetTableColumnsController.go", 23)
		}
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "#d41717", "ok", "None", 200)
		return c.Send(res)
	}
	if !checkGetTableColumnsRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "#d41717", "ok", "None", 200)
		return c.Send(res)
	}

	if !actions.CheckIfTableExists(obj.TableName) {
		res, _ := models.GetJSONResponse("table does not exist", "#d41717", "ok", "None", 200)
		return c.Send(res)
	}
	columns := actions.GetTableColumns(obj.Username, obj.Password, obj.Token, obj.TableName)

	// check response type
	if fmt.Sprintf("%T", columns) == "bool" {
		res, _ := models.GetJSONResponse("Error while fetching Array", "#d41717", "ok", "None", 200)
		return c.Send(res)
	}

	// if type of columns != bool
	return c.JSON(getTableColumnsResponse{
		Message: "successful",
		Alert:   "#1db004",
		Columns: columns,
	})
}

// checks the request
// struct fields cannot contain default value
func checkGetTableColumnsRequest(obj getTableColumnsRequest) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.TableName != ""
}
