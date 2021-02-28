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

func EditTableMinPermLvlController(c *fiber.Ctx) error {
	obj := new(editTableMinPermLvlRequest)
	err := c.BodyParser(obj)
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
	if !actions.MysqlLoginWithToken(obj.Username, obj.Password, obj.Token) {
		res, _ := models.GetJSONResponse("You do not have the permission to perform this", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	table := actions.GetTableByName(obj.TableName)
	conn := actions.GetConn()
	defer conn.Close()
	if !actions.CheckUserHasHigherPermission(conn, obj.Username, table.MinPermLvl, "") {
		res, _ := models.GetJSONResponse("You do not have the permission to perform this", "alert alert-warning", "ok", "None", 200)
		return c.Send(res)
	}
	actions.UpdateTableMinPermLvl(obj.TableName, obj.NewLvl)
	res, _ := models.GetJSONResponse("Successfully updated minimum permission level of table", "alert alert-success", "ok", "None", 200)
	return c.Send(res)
}

func checkEditTableMinPermLvlRequest(obj *editTableMinPermLvlRequest) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.TableName != "" && obj.NewLvl > 0
}
