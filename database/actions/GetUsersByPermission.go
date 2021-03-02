package actions

import "github.com/MathisBurger/OpenInventory-Backend/database/models"

//////////////////////////////////
// Queries all user with        //
// specific permission          //
//////////////////////////////////
func GetUsersByPermission(permission string) []models.UserModel {

	conn := GetConn()
	defer conn.Close()

	stmt, _ := conn.Prepare("SELECT * FROM `inv_users` WHERE `permissions` LIKE ?")
	defer stmt.Close()

	resp, _ := stmt.Query(permission)
	return models.UserModel{}.ParseAll(resp)
}
