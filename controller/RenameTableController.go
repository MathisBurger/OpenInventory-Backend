package controller

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/config"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

type renameTableRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Token     string `json:"token"`
	TableName string `json:"table_name"`
	NewName   string `json:"new_name"`
}

////////////////////////////////////////////////////////////////////
//                                                                //
//                     RenameTableController                      //
//               This controller renames the table                //
//            It requires renameTableRequest instance             //
//                                                                //
////////////////////////////////////////////////////////////////////
func RenameTableController(c *fiber.Ctx) error {

	// init and parse the request object
	obj := renameTableRequest{}
	err := json.Unmarshal(c.Body(), &obj)

	// check request
	if err != nil {
		if cfg, _ := config.ParseConfig(); cfg.ServerCFG.LogRequestErrors {
			utils.LogError(err.Error(), "RenameTableController.go", 24)
		}
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "#d41717", "ok", "None", 200)
		return c.Send(res)
	}
	if !checkRenameTableRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "#d41717", "ok", "None", 200)
		return c.Send(res)
	}

	// check login
	if !actions.MysqlLoginWithToken(obj.Username, obj.Password, obj.Token) {
		res, _ := models.GetJSONResponse("You do not have the permission to perform this", "#d41717", "ok", "None", 200)
		return c.Send(res)
	}

	conn := actions.GetConn()
	defer conn.Close()

	table := actions.GetTableByName(obj.TableName)

	// check permission
	if actions.CheckUserHasHigherPermission(conn, obj.Username, table.MinPermLvl, "") {

		if !actions.RenameTable(obj.TableName, obj.NewName) {
			res, _ := models.GetJSONResponse("Error while renaming table", "#d41717", "ok", "None", 200)
			return c.Send(res)
		}

		actions.UpdateTablename(obj.TableName, obj.NewName)

		res, _ := models.GetJSONResponse("Successfully updated tablename", "#1db004", "ok", "None", 200)
		return c.Send(res)
	}

	res, _ := models.GetJSONResponse("You do not have the permission to perform this", "#d41717", "ok", "None", 200)
	return c.Send(res)
}

// checks the request
// struct fields should not be default
func checkRenameTableRequest(obj renameTableRequest) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.TableName != "" && obj.NewName != ""
}
