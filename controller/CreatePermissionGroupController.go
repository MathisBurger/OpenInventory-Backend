package controller

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	OwnSQL "github.com/MathisBurger/OpenInventory-Backend/mysql"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type CreatePermissionGroupRequest struct {
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
	raw := string(c.Body())
	obj := CreatePermissionGroupRequest{}
	err := json.Unmarshal([]byte(raw), &obj)
	if err != nil {
		utils.LogError("[CreatePermissionGroupController.go, 28, InputError] " + err.Error())
		res, _ := models.GetJsonResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !checkCreatePermissionGroupRequest(obj) {
		res, _ := models.GetJsonResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !OwnSQL.MySQL_loginWithToken(obj.Username, obj.Password, obj.Token) {
		res, _ := models.GetJsonResponse("You do not have the permission to perform this", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	} else {
		permGroupInputStatus := checkPermissionGroupInput(obj)
		if permGroupInputStatus != nil {
			return c.Send(permGroupInputStatus)
		}
		conn := OwnSQL.GetConn()
		stmt, err := conn.Prepare("SELECT * FROM `inv_permissions` WHERE `name`=?")
		if err != nil {
			utils.LogError("[CreatePermissionGroupController.go, 43, SQL-StatementError] " + err.Error())
		}
		resp, err := stmt.Query("permission." + obj.PermissionInfo.Name)
		counter := 0
		for resp.Next() {
			counter += 1
		}
		defer resp.Close()
		if counter > 0 {
			res, _ := models.GetJsonResponse("This group already exists", "alert alert-warning", "ok", "None", 200)
			return c.Send(res)
		}
		stmt, err = conn.Prepare("INSERT INTO `inv_permissions` (`ID`, `name`, `color`, `permission-level`) VALUES (NULL, ?, ?, ?);")
		if err != nil {
			utils.LogError("[CreatePermissionGroupController.go, 57, SQL-StatementError] " + err.Error())
		}
		_, err = stmt.Exec("permission."+obj.PermissionInfo.Name, obj.PermissionInfo.ColorCode, obj.PermissionInfo.PermissionLevel)
		if err != nil {
			utils.LogError("[CreatePermissionGroupController.go, 61, SQL-StatementError] " + err.Error())
		}
		res, _ := models.GetJsonResponse("Created permissiongroup", "alert alert-success", "ok", "None", 200)
		return c.Send(res)
	}
}

func checkCreatePermissionGroupRequest(obj CreatePermissionGroupRequest) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.PermissionInfo.Name != "" && obj.PermissionInfo.ColorCode != "" && obj.PermissionInfo.PermissionLevel > 0
}

func checkPermissionGroupInput(obj CreatePermissionGroupRequest) []byte {
	if strings.Contains(obj.PermissionInfo.Name, ";") {
		resp, _ := models.GetJsonResponse("';' is not allowed in group name", "alert alert-danger", "ok", "None", 200)
		return resp
	}
	conn := OwnSQL.GetConn()
	if !OwnSQL.CheckUserHasHigherPermission(conn, obj.Username, obj.PermissionInfo.PermissionLevel, "") {
		defer conn.Close()
		resp, _ := models.GetJsonResponse("Your permission is not high enough", "alert alert-danger", "ok", "None", 200)
		return resp
	}
	defer conn.Close()
	return []byte{}
}
