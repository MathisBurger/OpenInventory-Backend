package controller

import (
	"github.com/MathisBurger/OpenInventory-Backend/config"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type createPermissionGroupRequest struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	Token          string `json:"token"`
	PermissionInfo struct {
		Name            string `json:"name"`
		ColorCode       string `json:"color_code"`
		PermissionLevel int    `json:"permission_level"`
	} `json:"permission_info"`
}

func CreatePermissionGroupController(c *fiber.Ctx) error {
	obj := new(createPermissionGroupRequest)
	err := c.BodyParser(obj)
	if err != nil {
		if cfg, _ := config.ParseConfig(); cfg.ServerCFG.LogRequestErrors {
			utils.LogError(err.Error(), "CreatePermissionGroupController.go", 29)
		}
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !checkCreatePermissionGroupRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !actions.MysqlLoginWithToken(obj.Username, obj.Password, obj.Token) {
		res, _ := models.GetJSONResponse("You do not have the permission to perform this", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	permGroupInputStatus := checkPermissionGroupInput(obj)
	if permGroupInputStatus != nil {
		return c.Send(permGroupInputStatus)
	}
	exists, _ := actions.GetPermissionByName(obj.PermissionInfo.Name)
	if exists {
		res, _ := models.GetJSONResponse("This group already exists", "alert alert-warning", "ok", "None", 200)
		return c.Send(res)
	}
	actions.InsertPermissionGroup(obj.PermissionInfo.Name, obj.PermissionInfo.ColorCode, obj.PermissionInfo.PermissionLevel)
	res, _ := models.GetJSONResponse("Created permissiongroup", "alert alert-success", "ok", "None", 200)
	return c.Send(res)
}

func checkCreatePermissionGroupRequest(obj *createPermissionGroupRequest) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.PermissionInfo.Name != "" && obj.PermissionInfo.ColorCode != "" && obj.PermissionInfo.PermissionLevel > 0
}

func checkPermissionGroupInput(obj *createPermissionGroupRequest) []byte {
	if strings.Contains(obj.PermissionInfo.Name, ";") {
		res, _ := models.GetJSONResponse("';' is not allowed in group name", "alert alert-danger", "ok", "None", 200)
		return res
	}
	conn := actions.GetConn()
	if !actions.CheckUserHasHigherPermission(conn, obj.Username, obj.PermissionInfo.PermissionLevel, "") {
		defer conn.Close()
		res, _ := models.GetJSONResponse("Your permission is not high enough", "alert alert-danger", "ok", "None", 200)
		return res
	}
	defer conn.Close()
	return nil
}
