package actions

import "github.com/MathisBurger/OpenInventory-Backend/database/models"

/////////////////////////////////
// Queries user by username    //
/////////////////////////////////
func GetUserByUsername(username string) (bool, models.UserModel) {

	conn := GetConn()
	defer conn.Close()

	stmt, _ := conn.Prepare("SELECT * FROM `inv_users` WHERE `username`=?")
	defer stmt.Close()

	resp, _ := stmt.Query(username)
	defer resp.Close()

	if !resp.Next() {
		return false, models.UserModel{}
	} else {
		return true, models.UserModel{}.Parse(resp)
	}
}
