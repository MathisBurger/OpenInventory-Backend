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

func MySQL_login(username string, password string) (bool, string) {
	conn := GetConn()
	hash := utils.HashWithSalt(password)
	stmt, err := conn.Prepare("SELECT * FROM inv_users WHERE displayname=? AND password=?")
	if err != nil {
		panic(err)
	}
	resp, err2 := stmt.Query(username, hash)
	if err2 != nil {
		panic(err2)
	}
	var answers []string
	for resp.Next() {
		var user UserStruct
		err2 = resp.Scan(&user.Username)
		if err2 != nil {
		}
		answers = append(answers, user.Username)
	}
	resp.Close()
	stmt.Close()
	conn.Close()
	if len(answers) == 1 {
		stmt, err = conn.Prepare("UPDATE inv_users SET token = ? WHERE displayname=?")
		if err != nil {
			panic(err)
		}
		token := utils.GenerateToken()
		_, _ = stmt.Exec(token, username)
		return true, token
	} else {
		return false, ""
	}

}

func MySQL_loginWithToken(username string, password string, token string) bool {
	conn := GetConn()
	hash := utils.HashWithSalt(password)
	stmt, err := conn.Prepare("SELECT * FROM inv_users WHERE displayname=? AND password=? AND token=?")
	if err != nil {
		panic(err)
	}
	resp, err2 := stmt.Query(username, hash, token)
	if err2 != nil {
		panic(err2)
	}
	var answers []string
	for resp.Next() {
		var user UserStruct
		err2 = resp.Scan(&user.Displayname)
		if err2 != nil {
		}
		answers = append(answers, user.Displayname)
	}
	resp.Close()
	stmt.Close()
	conn.Close()
	if len(answers) == 1 {
		return true
	} else {
		return false
	}

}

func MySQL_loginWithToken_ROOT(username string, password string, token string) bool {
	conn := GetConn()
	hash := utils.HashWithSalt(password)
	stmt, err := conn.Prepare("SELECT * FROM inv_users WHERE displayname=? AND password=? AND token=? AND root=1;")
	if err != nil {
		panic(err)
	}
	resp, err2 := stmt.Query(username, hash, token)
	if err2 != nil {
		panic(err2)
	}
	var answers []string
	for resp.Next() {
		var user UserStruct
		err2 = resp.Scan(&user.Displayname)
		if err2 != nil {
		}
		answers = append(answers, user.Displayname)
	}
	resp.Close()
	stmt.Close()
	conn.Close()
	if len(answers) == 1 {
		return true
	} else {
		return false
	}

}
