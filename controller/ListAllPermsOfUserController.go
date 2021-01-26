package controller

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	OwnSQL "github.com/MathisBurger/OpenInventory-Backend/mysql"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type ListAllPermsOfUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
	User     string `json:"user"`
}

type ListAllPermsOfUserResponse struct {
	Message     string                   `json:"message"`
	Permissions []models.PermissionModel `json:"permissions"`
	Status      string                   `json:"status"`
	HttpStatus  int                      `json:"http_status"`
}

func ListAllPermOfUserController(c *fiber.Ctx) error {
	raw := string(c.Body())
	obj := ListAllPermsOfUserRequest{}
	err := json.Unmarshal([]byte(raw), &obj)
	if err != nil {
		utils.LogError("[ListAllPermsOfUserController.go, 23, InputError] " + err.Error())
		res, _ := models.GetJsonResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !checkListAllPermsOfUserRequest(obj) {
		res, _ := models.GetJsonResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !OwnSQL.MySQL_loginWithToken(obj.Username, obj.Password, obj.Token) {
		res, _ := models.GetJsonResponse("You do not have the permission to perform this", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	} else {
		conn := OwnSQL.GetConn()
		stmt, err := conn.Prepare("SELECT `permissions` FROM `inv_users` WHERE `username`=?")
		if err != nil {
			utils.LogError("[ListAllPermsOfUserController.go, 38, SQL-StatementError] " + err.Error())
		}
		resp, err := stmt.Query(obj.User)
		if err != nil {
			utils.LogError("[ListAllPermsOfUserController.go, 42, SQL-StatementError] " + err.Error())
		}
		var perms string
		type cacheStruct struct {
			Permissions string `json:"permissions"`
		}
		for resp.Next() {
			var cache cacheStruct
			err = resp.Scan(&cache.Permissions)
			if err != nil {
				utils.LogError("[ListAllPermsOfUserController.go, 52, SQL-StatementError] " + err.Error())
			}
			perms = cache.Permissions
		}
		perm_names := strings.Split(perms, ";")
		var response []models.PermissionModel
		stmt, err = conn.Prepare("SELECT * FROM `inv_permissions` WHERE `name`=?")
		if err != nil {
			utils.LogError("[ListAllPermsOfUserController.go, 67, SQL-StatementError] " + err.Error())
		}
		for _, v := range perm_names {
			resp, err = stmt.Query(v)
			if err != nil {
				utils.LogError("[ListAllPermsOfUserController.go, 72, SQL-StatementError] " + err.Error())
			}
			for resp.Next() {
				var cache models.PermissionModel
				err = resp.Scan(&cache.ID, &cache.Name, &cache.Color, &cache.PermissionLevel)
				if err != nil {
					utils.LogError("[ListAllPermsOfUserController.go, 78, SQL-StatementError] " + err.Error())
				}
				response = append(response, cache)
			}
		}
		return c.JSON(ListAllPermsOfUserResponse{
			"Successfully fetched all user permissions",
			response,
			"ok",
			200,
		})
	}
}

func checkListAllPermsOfUserRequest(obj ListAllPermsOfUserRequest) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.User != ""
}
