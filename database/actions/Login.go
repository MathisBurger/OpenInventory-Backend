package actions

import (
	"github.com/MathisBurger/OpenInventory-Backend/database/models"
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

	var answers []models.UserModel
	for resp.Next() {
		var user models.UserModel
		_ = resp.Scan(&user.ID, &user.Username, &user.Password, &user.Token, &user.Permissions,
			&user.Root, &user.Mail, &user.Displayname, &user.RegisterDate, &user.Status)

		answers = append(answers, user)
	}

	if len(answers) == 1 {

		if utils.ValidateHash(password, answers[0].Password) {
			return true, ""
		}

		return false, ""
	}

	return false, ""
}
