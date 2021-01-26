package controller

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	OwnSQL "github.com/MathisBurger/OpenInventory-Backend/mysql"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func GetAllTablesController(c *fiber.Ctx) error {
	raw := string(c.Body())
	obj := models.LoginWithTokenRequest{}
	err := json.Unmarshal([]byte(raw), &obj)
	if err != nil {
		res, err := models.GetJsonResponse("Invaild JSON body", "alert alert-danger", "error", "None", 200)
		if err != nil {
			utils.LogError("[GetAllTablesController.go, 19, InputError] " + err.Error())
		}
		return c.Send(res)
	}
	if !checkCheckCredsRequestModel(obj) {
		res, _ := models.GetJsonResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !OwnSQL.MySQL_loginWithToken(obj.Username, obj.Password, obj.Token) {
		res, _ := models.GetJsonResponse("You do not have the permission to perform this", "alert alert-warning", "Failed", "None", 200)
		return c.Send(res)
	} else {
		tables := OwnSQL.GetAllTables(obj.Username, obj.Password, obj.Token)
		var compiledTables []string
		for _, table := range tables {
			cache := "['" + table.Name + "','" + strconv.Itoa(table.Entries) + "','" + table.CreatedAt + "']"
			compiledTables = append(compiledTables, cache)
		}
		msg := ""
		for _, str := range compiledTables {
			msg += str + ";"
		}
		mdl, err := models.GetJsonResponse(msg, "alert alert-success", "ok", "None", 200)
		if err != nil {
			utils.LogError("[GetAllTablesController.go, 39, ParsingError] " + err.Error())
		}
		return c.Send(mdl)
	}
}
