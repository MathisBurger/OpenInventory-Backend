package controller

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	OwnSQL "github.com/MathisBurger/OpenInventory-Backend/mysql"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

type ListAllPermissionGroupsResponse struct {
	Message          string                   `json:"message"`
	PermissionGroups []models.PermissionModel `json:"permission_groups"`
}

func ListAllPermissionGroupsController(c *fiber.Ctx) error {
	raw := string(c.Body())
	obj := models.LoginWithTokenRequest{}
	err := json.Unmarshal([]byte(raw), &obj)
	if err != nil {
		utils.LogError("[ListAllPermissionGroupsController.go, 16, InputError] " + err.Error())
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !checkCheckCredsRequestModel(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !OwnSQL.MysqlLoginWithToken(obj.Username, obj.Password, obj.Token) {
		res, _ := models.GetJSONResponse("You do not have the permission to perform this", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	conn := OwnSQL.GetConn()
	stmt, err := conn.Prepare("SELECT * FROM `inv_permissions`")
	if err != nil {
		utils.LogError("[ListAllPermissionGroupsController.go, 31, SQL-StatementError] " + err.Error())
	}
	resp, err := stmt.Query()
	if err != nil {
		utils.LogError("[ListAllPermissionGroupsController.go, 35, SQL-StatementError] " + err.Error())
	}
	var perms []models.PermissionModel
	for resp.Next() {
		var cache models.PermissionModel
		err = resp.Scan(&cache.ID, &cache.Name, &cache.Color, &cache.PermissionLevel)
		if err != nil {
			utils.LogError("[ListAllPermissionGroupsController.go, 42, SQL-StatementError] " + err.Error())
		}
		perms = append(perms, cache)
	}
	defer resp.Close()
	defer stmt.Close()
	defer conn.Close()
	return c.JSON(ListAllPermissionGroupsResponse{
		"Successfully fetched all permission groups",
		perms,
	})
}
