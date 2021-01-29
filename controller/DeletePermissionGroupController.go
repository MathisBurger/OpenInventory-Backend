package controller

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	OwnSQL "github.com/MathisBurger/OpenInventory-Backend/mysql"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type DeletePermissionGroupRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Token     string `json:"token"`
	GroupName string `json:"group_name"`
}

func DeletePermissionGroupController(c *fiber.Ctx) error {
	raw := string(c.Body())
	obj := DeletePermissionGroupRequest{}
	err := json.Unmarshal([]byte(raw), &obj)
	if err != nil {
		res, err := models.GetJsonResponse("Invaild JSON body", "alert alert-danger", "error", "None", 200)
		if err != nil {
			utils.LogError("[DeletePermissionGroupController.go, 25, InputError] " + err.Error())
		}
		return c.Send(res)
	}
	if !checkDeletePermissionGroupRequest(obj) {
		res, _ := models.GetJsonResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !OwnSQL.MySQL_loginWithToken(obj.Username, obj.Password, obj.Token) {
		res, _ := models.GetJsonResponse("You do not have the permission to perform this", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	} else {
		conn := OwnSQL.GetConn()
		if OwnSQL.CheckUserHasHigherPermission(conn, obj.Username, 0, "permission."+obj.GroupName) {
			stmt, err := conn.Prepare("SELECT `id`, `permissions` FROM `inv_users` WHERE `permissions` LIKE ?")
			if err != nil {
				utils.LogError("[DeletePermissionGroupController.go, 37, SQL-StatementError] " + err.Error())
			}
			type cacheStruct struct {
				ID          int    `json:"id"`
				Permissions string `json:"permissions"`
			}
			req := "%permission." + obj.GroupName + "%"
			resp, err := stmt.Query(req)
			if err != nil {
				utils.LogError("[DeletePermissionGroupController.go, 45, SQL-StatementError] " + err.Error())
			}
			var user []cacheStruct
			for resp.Next() {
				var cache cacheStruct
				err = resp.Scan(&cache.ID, &cache.Permissions)
				if err != nil {
					utils.LogError("[DeletePermissionGroupController.go, 52, SQL-StatementError] " + err.Error())
				}
				user = append(user, cache)
			}
			for _, val := range user {
				split := strings.Split(val.Permissions, ";")
				reduced := utils.RemoveValueFromArray(split, "permission."+obj.GroupName)
				editedPerms := reduced[0]
				for i, val2 := range reduced {
					if i == 0 {
						continue
					}
					editedPerms += ";" + val2
				}
				stmt, err = conn.Prepare("UPDATE `inv_users` SET `permissions`=? WHERE `id`=?")
				if err != nil {
					utils.LogError("[DeletePermissionGroupController.go, 68, SQL-StatementError] " + err.Error())
				}
				stmt.Exec(editedPerms, val.ID)
			}
			stmt, err = conn.Prepare("DELETE FROM `inv_permissions` WHERE `name`=?")
			if err != nil {
				utils.LogError("[DeletePermissionGroupController.go, 74, SQL-StatementError] " + err.Error())
			}
			_, err = stmt.Exec("permission." + obj.GroupName)
			if err != nil {
				utils.LogError("[DeletePermissionGroupController.go, 84, SQL-StatementError] " + err.Error())
			}
			defer resp.Close()
			defer stmt.Close()
			defer conn.Close()
			res, _ := models.GetJsonResponse("Successfully deleted PermissionGroup", "alert alert-success", "ok", "None", 200)
			return c.Send(res)
		} else {
			defer conn.Close()
			res, _ := models.GetJsonResponse("You do not have the permission to perform this", "alert alert-danger", "ok", "None", 200)
			return c.Send(res)
		}
	}
}

func checkDeletePermissionGroupRequest(obj DeletePermissionGroupRequest) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.GroupName != ""
}
