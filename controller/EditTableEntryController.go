package controller

import (
	"github.com/MathisBurger/OpenInventory-Backend/config"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

// ---------------------------------------------
//               editTableEntryRequest
// ---------------------------------------------
type editTableEntryRequest struct {
	Username  string                 `json:"username"`
	Password  string                 `json:"password"`
	Token     string                 `json:"token"`
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
	obj := new(editTableEntryRequest)
	err := c.BodyParser(obj)

	// check parsing error
	if err != nil {

		// log error if given in config
		if cfg, _ := config.ParseConfig(); cfg.ServerCFG.LogRequestErrors {
			utils.LogError(err.Error(), "EditTableEntryController.go", 23)
		}

		// return "Invalid JSON body" response
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}

	// check if request has been parsed correctly
	if !checkEditTableEntryRequest(obj) {

		// send "Wrong JSON syntax" response
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}

	// check login
	if actions.MysqlLoginWithToken(obj.Username, obj.Password, obj.Token) {

		table := actions.GetTableByName(obj.TableName)

		// get connection
		conn := actions.GetConn()
		defer conn.Close()

		// check higher permission
		if actions.CheckUserHasHigherPermission(conn, obj.Username, table.MinPermLvl, "") {

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

			// create values
			values = append(values, obj.ObjectID)

			// execute statement
			_, err = stmt.Exec(values...)

			// return failed response on error
			if err != nil {
				resp, _ := models.GetJSONResponse("Illegal row-map", "alert alert-danger", "ok", "None", 200)
				return c.Send(resp)
			}

			// send successful message
			res, _ := models.GetJSONResponse("Successfully updated entry", "alert alert-success", "ok", "None", 200)
			return c.Send(res)
		}

		// send permission error response
		res, _ := models.GetJSONResponse("You do not have the permission to perform this", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}

	// send permission error response
	res, _ := models.GetJSONResponse("You do not have the permission to perform this", "alert alert-danger", "ok", "None", 200)
	return c.Send(res)
}

/////////////////////////////////////////////////////////////
//                                                         //
//               checkEditTableEntryRequest                //
//      This function is checking the request object       //
//      It requires the editTableEntryRequest object       //
//                                                         //
/////////////////////////////////////////////////////////////
func checkEditTableEntryRequest(obj *editTableEntryRequest) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.TableName != "" && len(obj.Row) > 0 && obj.ObjectID > 0
}
