package actions

import "github.com/MathisBurger/OpenInventory-Backend/database/models"

func GetAllPermissions() []models.PermissionModel {
	conn := GetConn()
	defer conn.Close()
	stmt, _ := conn.Prepare("SELECT * FROM `inv_permissions`")
	defer stmt.Close()
	resp, _ := stmt.Query()
	defer resp.Close()
	return models.PermissionModel{}.ParseAll(resp)
}
