package controller

import (
	"github.com/MathisBurger/OpenInventory-Backend/config"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
	"strings"
)

// ---------------------------------------------
//            createTableRequest
//    This struct contains login credentials,
//         table-name, minPermLvl and row
// ---------------------------------------------
type createTableRequest struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	Token      string `json:"token"`
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

	// initializing the request object
	obj := new(createTableRequest)

	// parsing the body into the request object
	err := c.BodyParser(obj)

	// returns "Wrong JSON syntax" response if error is unequal nil
	if err != nil {

		// checks if request errors should be logged
		if cfg, _ := config.ParseConfig(); cfg.ServerCFG.LogRequestErrors {

			// log error
			utils.LogError(err.Error(), "CreateTableController.go", 26)
		}

		// returns response
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}

	// check if request has been parsed correctly
	if !checkCreateTableRequest(obj) {

		// returns "Wrong JSON syntax" response
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}

	// check table name length
	if !checkTableNameLength(obj.TableName) {

		// returns "Table name is too long" response
		res, _ := models.GetJSONResponse("Table name is too long", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}

	// create table || check status of table creation
	if actions.CreateTable(obj.Username, obj.Password, obj.Token, obj.TableName, parse(obj.RowConfig), obj.MinPermLvl) {

		// returns "successful" response if status true
		res, _ := models.GetJSONResponse("successful", "alert alert-success", "ok", "None", 200)
		return c.Send(res)
	}

	// returns "creation failed" response if status failed
	res, _ := models.GetJSONResponse("creation failed", "alert alert-danger", "ok", "None", 200)
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

/////////////////////////////////////////////////////////////
//                                                         //
//                  checkCreateTableRequest                //
//      This function is checking the request object       //
//       It requires the createTableRequest object         //
//                                                         //
/////////////////////////////////////////////////////////////
func checkCreateTableRequest(obj *createTableRequest) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.TableName != "" && obj.RowConfig != ""
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
