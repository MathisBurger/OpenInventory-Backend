package controller

import (
	"github.com/MathisBurger/OpenInventory-Backend/config"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

type removeTableEntryRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Token     string `json:"token"`
	TableName string `json:"table_name"`
	RowID     int    `json:"row_id"`
}

////////////////////////////////////////////////////////////////////
//                                                                //
//                  RemoveTableEntryController                    //
//       This controller fetches all permissions of user          //
//        It requires removeTableEntryRequest instance            //
//                                                                //
////////////////////////////////////////////////////////////////////
func RemoveTableEntryController(c *fiber.Ctx) error {

	// init and parse the request object
	obj := new(removeTableEntryRequest)
	err := c.BodyParser(obj)

	// check request
	if err != nil {
		if cfg, _ := config.ParseConfig(); cfg.ServerCFG.LogRequestErrors {
			utils.LogError(err.Error(), "RemoveTableEntryController.go", 25)
		}
		res, _ := models.GetJSONResponse("Invaild JSON body", "alert alert-danger", "error", "None", 200)
		return c.Send(res)
	}
	if !checkRemoveTableEntryRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}

	// check login
	if actions.MysqlLoginWithToken(obj.Username, obj.Password, obj.Token) {
		table := actions.GetTableByName(obj.TableName)
		conn := actions.GetConn()
		defer conn.Close()

		// check permission
		if actions.CheckUserHasHigherPermission(conn, obj.Username, table.MinPermLvl, "") {

			// check deletion status
			if !actions.DeleteTableEntry(obj.RowID, obj.TableName) {
				res, _ := models.GetJSONResponse("EntryID does not exist", "alert alert-warning", "ok", "None", 200)
				return c.Send(res)
			}

			actions.ChangeNumOfEntrysBy(obj.TableName, -1)

			res, _ := models.GetJSONResponse("Successfully deleted entry", "alert alert-success", "ok", "None", 200)
			return c.Send(res)
		}
		res, _ := models.GetJSONResponse("You do not have the permission perform this", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	res, _ := models.GetJSONResponse("You do not have the permission perform this", "alert alert-danger", "ok", "None", 200)
	return c.Send(res)
}

// checks the request
// struct fields should not be default
func checkRemoveTableEntryRequest(obj *removeTableEntryRequest) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.TableName != "" && obj.RowID > 0
}
