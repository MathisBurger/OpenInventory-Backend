package table_management

import (
	"encoding/json"
	"strings"

	"github.com/MathisBurger/OpenInventory-Backend/config"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/middleware"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

type createTableRequest struct {
	TableName  string `json:"table_name"`
	MinPermLvl int    `json:"min_perm_lvl"`
	RowConfig  string `json:"row_config"`
}

/////////////////////////////////////////////////////////////
//                                                         //
//                   CreateTableController                 //
//            This controller creates a new table          //
//          It requires login credentials and table        //
//                                                         //
/////////////////////////////////////////////////////////////
func CreateTableController(c *fiber.Ctx) error {

	// init and parse the request object
	obj := createTableRequest{}
	err := json.Unmarshal(c.Body(), &obj)

	// check request
	if err != nil {
		if cfg, _ := config.ParseConfig(); cfg.ServerCFG.LogRequestErrors {
			utils.LogError(err.Error(), "CreateTableController.go", 26)
		}

		res, _ := models.GetJSONResponse("Wrong JSON syntax", "#d41717", "ok", "None", 200)
		return c.Send(res)
	}
	if !checkCreateTableRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "#d41717", "ok", "None", 200)
		return c.Send(res)
	}

	// check table name length
	if !checkTableNameLength(obj.TableName) {
		res, _ := models.GetJSONResponse("Table name is too long", "#d41717", "ok", "None", 200)
		return c.Send(res)
	}

	if ok, ident := middleware.ValidateAccessToken(c); ok {

		actions.CreateTable(ident, obj.TableName, parse(obj.RowConfig), obj.MinPermLvl)
		res, _ := models.GetJSONResponse("successful", "#1db004", "ok", "None", 200)
		return c.Send(res)
	}

	res, _ := models.GetJSONResponse("creation failed", "#d41717", "ok", "None", 200)
	return c.Send(res)
}

// parsing string to []models.RowConfigModel
func parse(val string) (ans []models.RowConfigModel) {

	// split string with ',' as separator
	arr := strings.Split(val, ",")

	// iterate trough array
	for _, el := range arr {

		// remove bad chars
		raws := strings.ReplaceAll(el, "(", "")
		raws = strings.ReplaceAll(raws, ")", "")
		raws = strings.ReplaceAll(raws, " ", "_")
		raws = strings.ReplaceAll(raws, "[", "")
		raws = strings.ReplaceAll(raws, "]", "")
		raws = strings.ReplaceAll(raws, "-", "_")

		// split element with ';' as separator
		spl := strings.Split(raws, ";")

		// append split as RowConfigModel
		ans = append(ans, models.RowConfigModel{spl[0], checkTableName(spl[1])})
	}

	// return response
	return
}

// checks the request
// struct fields should not be default
func checkCreateTableRequest(obj createTableRequest) bool {
	return obj.TableName != "" && obj.RowConfig != ""
}

// replacing old with new chars
func checkTableName(name string) string {
	return strings.ReplaceAll(name, "-", "_")
}

// check if table name is too long or too short
// returns this as boolean
func checkTableNameLength(name string) bool {
	split := strings.Split(name, "")
	return len(split) < 16 && len(split) != 0
}
