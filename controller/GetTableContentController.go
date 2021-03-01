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
//       getTableContentRequest request
// ---------------------------------------------
type getTableContentRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Token     string `json:"token"`
	TableName string `json:"table_name"`
}

// ---------------------------------------------
//       getTableContentResponse response
// ---------------------------------------------
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
	obj := new(getTableContentRequest)
	err := c.BodyParser(obj)
	if err != nil {
		if cfg, _ := config.ParseConfig(); cfg.ServerCFG.LogRequestErrors {
			utils.LogError(err.Error(), "GetTableContentController.go", 24)
		}
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !checkGetTableContentRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !actions.MysqlLoginWithToken(obj.Username, obj.Password, obj.Token) {
		res, _ := models.GetJSONResponse("You do not have the permission to perform this", "alert alert-warning", "Failed", "None", 200)
		return c.Send(res)
	}
	stmt := "SELECT * FROM `table_" + obj.TableName + "`;"
	conn := actions.GetConn()
	defer conn.Close()
	json, err := utils.QueryToJson(conn, stmt)
	if err != nil {
		res, _ := models.GetJSONResponse("Invalid table name", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	return c.JSON(getTableContentResponse{
		Message:    "successful",
		Alert:      "alert alert-success",
		Status:     "ok",
		HttpStatus: 200,
		Elements:   strings.ReplaceAll(strings.ReplaceAll(string(json), "\n", ""), "\t", ""),
	})

}

// checks the request
// struct fields should not be default
func checkGetTableContentRequest(obj *getTableContentRequest) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.TableName != ""
}
