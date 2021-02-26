package controller

import (
	"github.com/MathisBurger/OpenInventory-Backend/config"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type createTableRequestModel struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	Token      string `json:"token"`
	TableName  string `json:"table_name"`
	MinPermLvl int    `json:"min_perm_lvl"`
	RowConfig  string `json:"row_config"`
}

func CreateTableController(c *fiber.Ctx) error {
	obj := new(createTableRequestModel)
	err := c.BodyParser(obj)
	if err != nil {
		if cfg, _ := config.ParseConfig(); cfg.ServerCFG.LogRequestErrors {
			utils.LogError(err.Error(), "CreateTableController.go", 26)
		}
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !checkCreateTableRequestModel(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !checkTableNameLength(obj.TableName) {
		res, _ := models.GetJSONResponse("Table name is too long", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	status := actions.CreateTable(obj.Username, obj.Password, obj.Token, obj.TableName, parse(obj.RowConfig), obj.MinPermLvl)
	if status {
		res, _ := models.GetJSONResponse("successful", "alert alert-success", "ok", "None", 200)
		return c.Send(res)
	}
	res, _ := models.GetJSONResponse("creation failed", "alert alert-danger", "ok", "None", 200)
	return c.Send(res)
}

func parse(val string) (ans []models.RowConfigModel) {
	arr := strings.Split(val, ",")
	for _, el := range arr {
		raws := strings.ReplaceAll(el, "(", "")
		raws = strings.ReplaceAll(raws, ")", "")
		raws = strings.ReplaceAll(raws, " ", "_")
		raws = strings.ReplaceAll(raws, "[", "")
		raws = strings.ReplaceAll(raws, "]", "")
		raws = strings.ReplaceAll(raws, "-", "_")
		spl := strings.Split(raws, ";")
		ans = append(ans, models.RowConfigModel{spl[0], checkTableName(spl[1])})
	}
	return
}

func checkTableName(name string) string {
	return strings.ReplaceAll(name, "-", "_")
}

func checkCreateTableRequestModel(obj *createTableRequestModel) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.TableName != "" && obj.RowConfig != ""
}

func checkTableNameLength(name string) bool {
	split := strings.Split(name, "")
	return len(split) < 16 && len(split) != 0
}
