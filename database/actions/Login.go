package actions

import (
	"github.com/MathisBurger/OpenInventory-Backend/utils"
)

////////////////////////////////////////
// Checks login status of user        //
// by username and password           //
////////////////////////////////////////
func MysqlLogin(username string, password string) (bool, string) {

	conn := GetConn()
	defer conn.Close()

	stmt, _ := conn.Prepare("SELECT * FROM inv_users WHERE displayname=?")
	defer stmt.Close()

	resp, _ := stmt.Query(username)
	defer resp.Close()

	if exists, usr := GetUserByUsername(username); exists {

		if utils.ValidateHash(password, usr.Password) {
			return true, ""
		}

		return false, ""
	}

	return false, ""
}
