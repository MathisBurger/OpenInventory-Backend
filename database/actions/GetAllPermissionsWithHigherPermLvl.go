package actions

import "github.com/MathisBurger/OpenInventory-Backend/database/models"

/////////////////////////////////////////////////
// Queries all permission groups with higher   //
// permission level than the given             //
/////////////////////////////////////////////////
func GetAllPermissionsWithHigherPermLvl(minPermLvl int) []models.PermissionModel {

	conn := GetConn()
	defer conn.Close()

	stmt, _ := conn.Prepare("SELECT * FROM `inv_permissions` WHERE `permission-level`>=?")
	defer stmt.Close()

	resp, _ := stmt.Query(minPermLvl)
	defer resp.Close()

	return models.PermissionModel{}.ParseAll(resp)
}
