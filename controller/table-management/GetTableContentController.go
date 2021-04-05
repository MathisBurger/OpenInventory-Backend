package table_management

import (
	"strings"

	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/middleware"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

type getTableContentRequest struct {
	TableName string `json:"table_name"`
}

type getTableContentResponse struct {
	Message    string `json:"message"`
	Alert      string `json:"alert"`
	Status     string `json:"status"`
	HttpStatus int    `json:"httpStatus"`
	Elements   string `json:"elements"`
}

////////////////////////////////////////////////////////////////////
//                                                                //
//                   GetTableContentController                    //
//         This controller fetches content of given table         //
//          It requires getTableContentRequest instance           //
//                                                                //
////////////////////////////////////////////////////////////////////
func GetTableContentController(c *fiber.Ctx) error {

	// init and parse the request object
	obj := getTableContentRequest{
		TableName: c.Query("table_name", ""),
	}

	// check request
	if !checkGetTableContentRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "#d41717", "ok", "None", 200)
		return c.Send(res)
	}

	// check login
	if ok, ident := middleware.ValidateAccessToken(c); ok {

		tbl := actions.GetTableByName(obj.TableName).MinPermLvl

		conn := actions.GetConn()
		defer conn.Close()

		// check if permission of user is high enough
		if !actions.CheckUserHasHigherPermission(conn, ident, tbl, "") {
			res, _ := models.GetJSONResponse("Your permission is not high enough", "#d41717", "ok", "None", 200)
			return c.Send(res)
		}

		// query table as json
		stmt := "SELECT * FROM `table_" + obj.TableName + "`;"
		json, err := utils.QueryToJson(conn, stmt)

		if err != nil {
			res, _ := models.GetJSONResponse("Invalid table name", "#d41717", "ok", "None", 200)
			return c.Send(res)
		}

		return c.JSON(getTableContentResponse{
			Message:    "successful",
			Alert:      "#1db004",
			Status:     "ok",
			HttpStatus: 200,
			Elements:   strings.ReplaceAll(strings.ReplaceAll(string(json), "\n", ""), "\t", ""),
		})
	}

	res, _ := models.GetJSONResponse("You do not have the permission to perform this", "alert alert-warning", "Failed", "None", 200)
	return c.Send(res)
}

// checks the request
// struct fields should not be default
func checkGetTableContentRequest(obj getTableContentRequest) bool {
	return obj.TableName != ""
}
