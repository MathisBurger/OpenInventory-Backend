package table_management

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/config"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type renameTableColumnRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Token     string `json:"token"`
	TableName string `json:"table_name"`
	OldName   string `json:"old_name"`
	NewName   string `json:"new_name"`
}

////////////////////////////////////////////////////////////////////
//                                                                //
//                  RenameTableColumnController                   //
//             This controller renames a table column             //
//         It requires renameTableColumnRequest instance          //
//                                                                //
////////////////////////////////////////////////////////////////////
func RenameTableColumnController(c *fiber.Ctx) error {

	// init and parse the request object
	obj := renameTableColumnRequest{}
	err := json.Unmarshal(c.Body(), &obj)

	// check request
	if err != nil {
		if cfg, _ := config.ParseConfig(); cfg.ServerCFG.LogRequestErrors {
			utils.LogError(err.Error(), "RenameTableColumnController.go", 26)
		}
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "#d41717", "ok", "None", 200)
		return c.Send(res)
	}
	if !checkRenameTableColumnRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "#d41717", "ok", "None", 200)
		return c.Send(res)
	}

	// check login
	if !actions.MysqlLoginWithToken(obj.Username, obj.Password, obj.Token) {
		res, _ := models.GetJSONResponse("You do not have the permission to perform this", "#d41717", "Failed", "None", 200)
		return c.Send(res)
	}

	conn := actions.GetConn()
	defer conn.Close()

	columns := actions.GetTableColumns(obj.Username, obj.Password, obj.Token, obj.TableName)

	// check if user has permission for table
	if len(columns) == 0 {
		resp, _ := models.GetJSONResponse("You do not have the permission to perform this", "#d41717", "ok", "None", 200)
		return c.Send(resp)
	}

	table := actions.GetTableByName(obj.TableName)

	// check permission
	if actions.CheckUserHasHigherPermission(conn, obj.Username, table.MinPermLvl, "") {

		// iterate trough all columns
		for _, val := range columns {

			// check if val == column to rename
			if val.COLUMN_NAME == obj.OldName {

				// get params from column
				var length string
				if val.MAX_LENGTH == nil {
					length = ""
				} else {
					i, _ := val.MAX_LENGTH.(int64)
					length = strconv.Itoa(int(i))
				}
				if val.DATA_TYPE == "int" {
					length = "11"
				}

				// rename column
				if !actions.RenameTableColumn(obj.TableName, obj.OldName, obj.NewName, val.DATA_TYPE, length) {
					res, _ := models.GetJSONResponse("Error while changing column name", "#d41717", "ok", "None", 200)
					return c.Send(res)
				}

				res, _ := models.GetJSONResponse("Successfully changed column name", "#1db004", "ok", "None", 200)
				return c.Send(res)
			}
		}
		res, _ := models.GetJSONResponse("Column not found", "alert alert-warning", "ok", "None", 200)
		return c.Send(res)
	}
	res, _ := models.GetJSONResponse("You do not have the permission to do this", "alert alert-warning", "ok", "None", 200)
	return c.Send(res)
}

// checks the request
// struct fields should not be default
func checkRenameTableColumnRequest(obj renameTableColumnRequest) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.TableName != "" && obj.NewName != "" && obj.OldName != ""
}
