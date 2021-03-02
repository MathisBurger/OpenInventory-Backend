package actions

import (
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"time"
)

// cache struct
type UserStruct struct {
	Id           int       `json:"id"`
	Username     string    `json:"username"`
	Password     string    `json:"password"`
	Token        string    `json:"token"`
	Root         bool      `json:"root"`
	Mail         string    `json:"mail"`
	Displayname  string    `json:"displayname"`
	RegisterDate time.Time `json:"register_date"`
	Status       string    `json:"status"`
}

// cache struct
type DisplayNameStruct struct {
	Displayname string `json:"displayname"`
}

////////////////////////////////////////
// Checks login status of user        //
// by username and password           //
////////////////////////////////////////
func MysqlLogin(username string, password string) (bool, string) {

	conn := GetConn()
	defer conn.Close()

	hash := utils.HashWithSalt(password)

	stmt, _ := conn.Prepare("SELECT * FROM inv_users WHERE displayname=? AND password=?")
	defer stmt.Close()

	resp, _ := stmt.Query(username, hash)
	defer resp.Close()

	var answers []string
	for resp.Next() {
		var user UserStruct
		_ = resp.Scan(&user.Username)

		answers = append(answers, user.Username)
	}

	if len(answers) == 1 {
		stmt, _ = conn.Prepare("UPDATE inv_users SET token = ? WHERE displayname=?")
		defer stmt.Close()

		token := utils.GenerateToken()

		_, _ = stmt.Exec(token, username)
		return true, token
	}
	return false, ""

}

//////////////////////////////////////
// Checks login status by           //
// username, password and token     //
//////////////////////////////////////
func MysqlLoginWithToken(username string, password string, token string) bool {

	conn := GetConn()
	defer conn.Close()

	hash := utils.HashWithSalt(password)

	stmt, _ := conn.Prepare("SELECT `displayname` FROM inv_users WHERE displayname=? AND password=? AND token=?")
	defer stmt.Close()

	resp, _ := stmt.Query(username, hash, token)
	defer stmt.Close()

	var answers []string
	for resp.Next() {
		var user DisplayNameStruct
		_ = resp.Scan(&user.Displayname)

		answers = append(answers, user.Displayname)
	}

	return len(answers) == 1
}

//////////////////////////////////////////
// Checks login as root status by       //
// username, password and token         //
//////////////////////////////////////////
func MysqlLoginWithTokenRoot(username string, password string, token string) bool {

	conn := GetConn()
	defer conn.Close()

	hash := utils.HashWithSalt(password)

	stmt, _ := conn.Prepare("SELECT * FROM inv_users WHERE displayname=? AND password=? AND token=? AND root=1;")
	defer stmt.Close()

	resp, _ := stmt.Query(username, hash, token)
	defer resp.Close()

	var answers []string
	for resp.Next() {
		var user UserStruct
		_ = resp.Scan(&user.Displayname)

		answers = append(answers, user.Displayname)
	}

	return len(answers) == 1
}
