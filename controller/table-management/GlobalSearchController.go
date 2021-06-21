package table_management

import (
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/middleware"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type GlobalSearchResponse struct {
	Data []string `json:"data"`
}

////////////////////////////////////////////////////////////////////
//                                                                //
//                   GetTableContentController                    //
//         This controller fetches content of given table         //
//          It requires getTableContentRequest instance           //
//                                                                //
////////////////////////////////////////////////////////////////////
func GlobalSearchController(c *fiber.Ctx) error {

	if ok, ident := middleware.ValidateAccessToken(c); ok {
		_, user := actions.GetUserByUsername(ident)
		tables := actions.GetAllTables(user.Displayname)
		var content []string
		conn := actions.GetConn()
		defer conn.Close()

		for _, table := range tables {
			stmt := "SELECT * FROM `table_" + table.Name + "`;"
			json, _ := utils.QueryToJson(conn, stmt)
			content = append(content, strings.ReplaceAll(strings.ReplaceAll(string(json), "\n", ""), "\t", ""))
		}
		return c.JSON(GlobalSearchResponse{
			Data: content,
		})
	} else {
		res, _ := models.GetJSONResponse("You do not have the permission to perform this", "alert alert-warning", "Failed", "None", 200)
		return c.Send(res)
	}

}
