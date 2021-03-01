package controller

import (
	"github.com/MathisBurger/OpenInventory-Backend/config"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

type deleteTableRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Token     string `json:"token"`
	TableName string `json:"table_name"`
}

/////////////////////////////////////////////////////////////
//                                                         //
//                 DeleteTableController                   //
//            This controller deletes an table             //
//         It requires deleteTableRequest instance         //
//                                                         //
/////////////////////////////////////////////////////////////
func DeleteTableController(c *fiber.Ctx) error {

	// init and parse the request object
	obj := new(deleteTableRequest)
	err := c.BodyParser(obj)

	// check request
	if err != nil {
		if cfg, _ := config.ParseConfig(); cfg.ServerCFG.LogRequestErrors {
			utils.LogError(err.Error(), "DeleteTableController.go", 23)
		}

		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !checkDeleteTableRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}

	// check login status
	if actions.MysqlLoginWithToken(obj.Username, obj.Password, obj.Token) {

		conn := actions.GetConn()
		defer conn.Close()

		table := actions.GetTableByName(obj.TableName)

		// check permission of user
		if actions.CheckUserHasHigherPermission(conn, obj.Username, table.MinPermLvl, "") {

			actions.DropTable(obj.TableName)

			res, _ := models.GetJSONResponse("Successfully deleted table", "alert alert-success", "ok", "None", 200)
			return c.Send(res)
		}

		res, _ := models.GetJSONResponse("You do not have the permission to perform this", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)

	}

	res, _ := models.GetJSONResponse("You do not have the permission to perform this", "alert alert-danger", "ok", "None", 200)
	return c.Send(res)

}

// checks the request
// struct fields should not be default
func checkDeleteTableRequest(obj *deleteTableRequest) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.TableName != ""
}
