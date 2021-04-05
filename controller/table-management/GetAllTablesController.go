package table_management

import (
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/middleware"
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

	// check login
	if ok, ident := middleware.ValidateAccessToken(c); ok {

		tables := actions.GetAllTables(ident)

		return c.JSON(getAllTablesResponse{
			"Successfully queried all tables",
			"#1db004",
			tables,
		})
	}

	res, _ := models.GetJSONResponse("You do not have the permission to perform this", "alert alert-warning", "Failed", "None", 200)
	return c.Send(res)
}
