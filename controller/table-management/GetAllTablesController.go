package table_management

import (
	"github.com/MathisBurger/OpenInventory-Backend/controller/general"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/models"
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
	obj := models.LoginWithTokenRequest{
		Username: c.Query("username", ""),
		Password: c.Query("password", ""),
		Token:    c.Query("token", ""),
	}

	// check request
	if !general.CheckCheckCredsRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "#d41717", "ok", "None", 200)
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
		"#1db004",
		tables,
	})
}
