package controller

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	OwnSQL "github.com/MathisBurger/OpenInventory-Backend/mysql"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

type ListAllPermGroupsOfTableRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Token     string `json:"token"`
	TableName string `json:"table_name"`
}

type ListAllPermGroupsOfTableResponse struct {
	PermGroups []models.PermissionModel `json:"perm_groups"`
	Message    string                   `json:"message"`
}

func ListAllPermGroupsOfTableController(c *fiber.Ctx) error {
	raw := string(c.Body())
	obj := ListAllPermGroupsOfTableRequest{}
	err := json.Unmarshal([]byte(raw), &obj)
	if err != nil {
		utils.LogError("[ListAllPermGroupsOfTableController.go, 23, InputError] " + err.Error())
		res, _ := models.GetJsonResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !checkListAllPermGroupsOfTableRequest(obj) {
		res, _ := models.GetJsonResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !OwnSQL.MySQL_loginWithToken(obj.Username, obj.Password, obj.Token) {
		res, _ := models.GetJsonResponse("You do not have the permission to perform this", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	} else {
		conn := OwnSQL.GetConn()
		stmt, err := conn.Prepare("SELECT `min-perm-lvl` FROM `inv_tables` WHERE `name`=?")
		if err != nil {
			utils.LogError("[ListAllPermGroupsOfTableController.go, 38, SQL-StatementError] " + err.Error())
		}
		type cacheStruct struct {
			MinPermlvl int `json:"min-perm-lvl"`
		}
		var min_perm_lvl int
		resp, err := stmt.Query(obj.TableName)
		if err != nil {
			utils.LogError("[ListAllPermGroupsOfTableController.go, 38, SQL-StatementError] " + err.Error())
		}
		for resp.Next() {
			var cache cacheStruct
			err = resp.Scan(&cache.MinPermlvl)
			if err != nil {
				utils.LogError("[ListAllPermGroupsOfTableController.go, 52, SQL-StatementError] " + err.Error())
			}
			min_perm_lvl = cache.MinPermlvl
		}
		if !OwnSQL.CheckUserHasHigherPermission(conn, obj.Username, min_perm_lvl, "") {
			defer resp.Close()
			defer stmt.Close()
			defer conn.Close()
			res, _ := models.GetJsonResponse("Your permission is not high enough to view this table", "alert alert-danger", "ok", "None", 200)
			return c.Send(res)
		} else {
			stmt, err = conn.Prepare("SELECT * FROM `inv_permissions` WHERE `permission-level`>=?")
			if err != nil {
				utils.LogError("[ListAllPermGroupsOfTableController.go, 62, SQL-StatementError] " + err.Error())
			}
			resp, err = stmt.Query(min_perm_lvl)
			if err != nil {
				utils.LogError("[ListAllPermGroupsOfTableController.go, 66, SQL-StatementError] " + err.Error())
			}
			var response []models.PermissionModel
			for resp.Next() {
				var cache models.PermissionModel
				err = resp.Scan(&cache.ID, &cache.Name, &cache.Color, &cache.PermissionLevel)
				if err != nil {
					utils.LogError("[ListAllPermsOfUserController.go, 78, SQL-StatementError] " + err.Error())
				}
				response = append(response, cache)
			}
			defer resp.Close()
			defer stmt.Close()
			defer conn.Close()
			return c.JSON(ListAllPermGroupsOfTableResponse{
				response,
				"Successfully fetched all permissiongroups of table",
			})
		}
	}
}

func checkListAllPermGroupsOfTableRequest(obj ListAllPermGroupsOfTableRequest) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.TableName != ""
}
