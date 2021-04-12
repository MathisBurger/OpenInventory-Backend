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

type editTableEntryRequest struct {
	TableName string                 `json:"table_name"`
	ObjectID  int                    `json:"object_id"`
	Row       map[string]interface{} `json:"row"`
}

/////////////////////////////////////////////////////////////
//                                                         //
//                EditTableEntryController                 //
//    This controller changes values of the given row      //
//       It requires editTableEntryRequest instance        //
//                                                         //
/////////////////////////////////////////////////////////////
func EditTableEntryController(c *fiber.Ctx) error {

	// init and parse the request object
	obj := editTableEntryRequest{}
	decrypted, err := e2e.DecryptBytes(c.Body())
	if err != nil {
		return c.SendStatus(400)
	}
	err = json.Unmarshal(decrypted, &obj)

	// check request
	if err != nil {
		if cfg, _ := config.ParseConfig(); cfg.ServerCFG.LogRequestErrors {
			utils.LogError(err.Error(), "EditTableEntryController.go", 23)
		}
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "#d41717", "ok", "None", 200)
		return c.Send(res)
	}
	if !checkEditTableEntryRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "#d41717", "ok", "None", 200)
		return c.Send(res)
	}

	// check login
	if ok, ident := middleware.ValidateAccessToken(c); ok {

		table := actions.GetTableByName(obj.TableName)

		conn := actions.GetConn()
		defer conn.Close()

		// check higher permission
		if actions.CheckUserHasHigherPermission(conn, ident, table.MinPermLvl, "") {

			// build sql statement for editing table entry
			sql := "UPDATE `table_" + obj.TableName + "` SET "
			first_completed := false
			var values []interface{}
			for k, v := range obj.Row {
				if k != "id" {
					if !first_completed {
						sql += "`" + k + "`=?"
						values = append(values, v)
						first_completed = true
					} else {
						sql += ", `" + k + "`=?"
						values = append(values, v)
					}
				}
			}
			sql += " WHERE `id`=?"

			// prepare statement
			stmt, err := conn.Prepare(sql)
			defer stmt.Close()
			if err != nil {
				utils.LogError(err.Error(), "EditTableEntryController.go", 56)
			}

			values = append(values, obj.ObjectID)

			_, err = stmt.Exec(values...)

			if err != nil {
				resp, _ := models.GetJSONResponse("Illegal row-map", "#d41717", "ok", "None", 200)
				return c.Send(resp)
			}

			res, _ := models.GetJSONResponse("Successfully updated entry", "#1db004", "ok", "None", 200)
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
func checkEditTableEntryRequest(obj editTableEntryRequest) bool {
	return obj.TableName != "" && len(obj.Row) > 0 && obj.ObjectID > 0
}
