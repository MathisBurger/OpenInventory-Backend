package actions

import (
	"github.com/MathisBurger/OpenInventory-Backend/database/models"
	"strings"
)

////////////////////////////////////////
// Queries all permissions of user    //
////////////////////////////////////////
func GetPermissionsOfUser(username string) []models.PermissionModel {

	_, user := GetUserByUsername(username)
	perms := strings.Split(user.Permissions, ";")

	conn := GetConn()
	defer conn.Close()

	stmt, _ := conn.Prepare("SELECT * FROM `inv_permissions` WHERE `name`=?")
	defer stmt.Close()

	var response []models.PermissionModel

	// Appends all permissions to array
	for _, v := range perms {
		resp, _ := stmt.Query(v)
		defer resp.Close()
		for resp.Next() {
			var cache models.PermissionModel
			_ = resp.Scan(&cache.ID, &cache.Name, &cache.Color, &cache.PermissionLevel)
			response = append(response, cache)
		}
	}

	return response
}
