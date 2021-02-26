package actions

import "github.com/MathisBurger/OpenInventory-Backend/database/models"

func GetPermissionByName(name string) (bool, models.PermissionModel) {
	conn := GetConn()
	defer conn.Close()
	stmt, _ := conn.Prepare("SELECT * FROM `inv_permissions` WHERE `name`=?;")
	defer stmt.Close()
	resp, _ := stmt.Query(name)
	defer resp.Close()
	if !resp.Next() {
		return false, models.PermissionModel{}
	} else {
		return true, models.PermissionModel{}.Parse(resp)
	}
}
