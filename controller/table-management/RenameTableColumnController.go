package table_management

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/e2e"
	"strconv"

	"github.com/MathisBurger/OpenInventory-Backend/config"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/middleware"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

type renameTableColumnRequest struct {
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
	decrypted, err := e2e.DecryptBytes(c.Body())
	if err != nil {
		return c.SendStatus(400)
	}
	err = json.Unmarshal(decrypted, &obj)

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
	if ok, ident := middleware.ValidateAccessToken(c); ok {
		conn := actions.GetConn()
		defer conn.Close()

		columns := actions.GetTableColumns(ident, obj.TableName)

		// check if user has permission for table
		if len(columns) == 0 {
			resp, _ := models.GetJSONResponse("You do not have the permission to perform this", "#d41717", "ok", "None", 200)
			return c.Send(resp)
		}

		table := actions.GetTableByName(obj.TableName)

		// check permission
		if actions.CheckUserHasHigherPermission(conn, ident, table.MinPermLvl, "") {

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

	res, _ := models.GetJSONResponse("You do not have the permission to perform this", "#d41717", "Failed", "None", 200)
	return c.Send(res)
}

// checks the request
// struct fields should not be default
func checkRenameTableColumnRequest(obj renameTableColumnRequest) bool {
	return obj.TableName != "" && obj.NewName != "" && obj.OldName != ""
}
