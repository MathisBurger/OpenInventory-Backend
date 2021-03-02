package actions

import "github.com/MathisBurger/OpenInventory-Backend/database/models"

//////////////////////////////
// Queries all user         //
//////////////////////////////
func GetAllUser() []models.UserModel {

	conn := GetConn()
	defer conn.Close()

	stmt, _ := conn.Prepare("SELECT * FROM `inv_users`")
	defer stmt.Close()

	resp, _ := stmt.Query()
	defer resp.Close()

	return models.UserModel{}.ParseAll(resp)
}
