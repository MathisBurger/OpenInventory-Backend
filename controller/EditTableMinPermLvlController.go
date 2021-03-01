package controller

import (
	"github.com/MathisBurger/OpenInventory-Backend/config"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

type editTableMinPermLvlRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Token     string `json:"token"`
	TableName string `json:"table_name"`
	NewLvl    int    `json:"new_lvl"`
}

////////////////////////////////////////////////////////////////////
//                                                                //
//                EditTableMinPermLvlController                   //
//    This controller changes the minPermLvl of the given table   //
//          It requires editTableMinPermLvlRequest instance       //
//                                                                //
////////////////////////////////////////////////////////////////////
func EditTableMinPermLvlController(c *fiber.Ctx) error {

	// init and parse the request object
	obj := new(editTableMinPermLvlRequest)
	err := c.BodyParser(obj)

	// check request
	if err != nil {
		if cfg, _ := config.ParseConfig(); cfg.ServerCFG.LogRequestErrors {
			utils.LogError(err.Error(), "EditTableMinPermLvlController.go", 24)
		}
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !checkEditTableMinPermLvlRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}

	// check login
	if !actions.MysqlLoginWithToken(obj.Username, obj.Password, obj.Token) {
		res, _ := models.GetJSONResponse("You do not have the permission to perform this", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}

	table := actions.GetTableByName(obj.TableName)
	conn := actions.GetConn()
	defer conn.Close()

	// check higher permission
	if !actions.CheckUserHasHigherPermission(conn, obj.Username, table.MinPermLvl, "") {
		res, _ := models.GetJSONResponse("You do not have the permission to perform this", "alert alert-warning", "ok", "None", 200)
		return c.Send(res)
	}

	// update min perm lvl
	actions.UpdateTableMinPermLvl(obj.TableName, obj.NewLvl)

	res, _ := models.GetJSONResponse("Successfully updated minimum permission level of table", "alert alert-success", "ok", "None", 200)
	return c.Send(res)
}

// checks the request
// struct fields should not be default
func checkEditTableMinPermLvlRequest(obj *editTableMinPermLvlRequest) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.TableName != "" && obj.NewLvl > 0
}
