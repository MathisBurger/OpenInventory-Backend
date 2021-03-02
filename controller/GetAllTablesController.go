package controller

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/config"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

type getAllTablesResponse struct {
	Message string              `json:"message"`
	Alert   string              `json:"alert"`
	Tables  []models.TableModel `json:"tables"`
}

////////////////////////////////////////////////////////////////////
//                                                                //
//                     GetAllTablesController                     //
//               This controller fetches all tables               //
//        It requires models.LoginWithTokenRequest instance       //
//                                                                //
////////////////////////////////////////////////////////////////////
func GetAllTablesController(c *fiber.Ctx) error {

	// init and parse the request object
	obj := models.LoginWithTokenRequest{}
	err := json.Unmarshal(c.Body(), &obj)

	// check request
	if err != nil {
		if cfg, _ := config.ParseConfig(); cfg.ServerCFG.LogRequestErrors {
			utils.LogError(err.Error(), "GetAllTablesController.go", 17)
		}
		res, _ := models.GetJSONResponse("Invaild JSON body", "alert alert-danger", "error", "None", 200)
		return c.Send(res)
	}
	if !checkCheckCredsRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}

	// check login
	if !actions.MysqlLoginWithToken(obj.Username, obj.Password, obj.Token) {
		res, _ := models.GetJSONResponse("You do not have the permission to perform this", "alert alert-warning", "Failed", "None", 200)
		return c.Send(res)
	}

	tables := actions.GetAllTables(obj.Username, obj.Password, obj.Token)

	return c.JSON(getAllTablesResponse{
		"Successfully queried all tables",
		"alert alert-success",
		tables,
	})
}
