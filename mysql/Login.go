package OwnSQL

import (
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"time"
)

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

type DisplayNameStruct struct {
	Displayname string `json:"displayname"`
}

func MysqlLogin(username string, password string) (bool, string) {
	conn := GetConn()
	hash := utils.HashWithSalt(password)
	stmt, err := conn.Prepare("SELECT * FROM inv_users WHERE displayname=? AND password=?")
	if err != nil {
		utils.LogError("[Login.go, 27, SQL-StatementError] " + err.Error())
	}
	resp, err := stmt.Query(username, hash)
	if err != nil {
		utils.LogError("[Login.go, 31, SQL-StatementError] " + err.Error())
	}
	var answers []string
	for resp.Next() {
		var user UserStruct
		err = resp.Scan(&user.Username)
		if err != nil {
			utils.LogError("[Login.go, 31, SQL-ScanningError] " + err.Error())
		}
		answers = append(answers, user.Username)
	}
	defer resp.Close()
	defer stmt.Close()
	defer conn.Close()
	if len(answers) == 1 {
		stmt, err = conn.Prepare("UPDATE inv_users SET token = ? WHERE displayname=?")
		if err != nil {
			utils.LogError("[Login.go, 48, SQL-StatementError] " + err.Error())
		}
		token := utils.GenerateToken()
		_, _ = stmt.Exec(token, username)
		return true, token
	}
	return false, ""

}

func MysqlLoginWithToken(username string, password string, token string) bool {
	conn := GetConn()
	hash := utils.HashWithSalt(password)
	stmt, err := conn.Prepare("SELECT `displayname` FROM inv_users WHERE displayname=? AND password=? AND token=?")
	if err != nil {
		utils.LogError("[Login.go, 64, SQL-StatementError] " + err.Error())
	}
	resp, err := stmt.Query(username, hash, token)
	if err != nil {
		utils.LogError("[Login.go, 68, SQL-StatementError] " + err.Error())
	}
	var answers []string
	for resp.Next() {
		var user DisplayNameStruct
		err = resp.Scan(&user.Displayname)
		if err != nil {
			utils.LogError("[Login.go, 75, SQL-ScanningError] " + err.Error())
		}
		answers = append(answers, user.Displayname)
	}
	defer resp.Close()
	defer stmt.Close()
	defer conn.Close()
	if len(answers) == 1 {
		return true
	} else {
		return false
	}

}

func MysqlLoginWithTokenRoot(username string, password string, token string) bool {
	conn := GetConn()
	hash := utils.HashWithSalt(password)
	stmt, err := conn.Prepare("SELECT * FROM inv_users WHERE displayname=? AND password=? AND token=? AND root=1;")
	if err != nil {
		utils.LogError("[Login.go, 95, SQL-StatementError] " + err.Error())
	}
	resp, err := stmt.Query(username, hash, token)
	if err != nil {
		utils.LogError("[Login.go, 99, SQL-StatementError] " + err.Error())
	}
	var answers []string
	for resp.Next() {
		var user UserStruct
		err = resp.Scan(&user.Displayname)
		if err != nil {
			utils.LogError("[Login.go, 106, SQL-ScanningError] " + err.Error())
		}
		answers = append(answers, user.Displayname)
	}
	defer resp.Close()
	defer stmt.Close()
	defer conn.Close()
	if len(answers) == 1 {
		return true
	} else {
		return false
	}

}
