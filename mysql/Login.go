package OwnSQL

import (
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"time"
)

type UserStruct struct {
	id            int       `json:"id"`
	username      string    `json:"username"`
	password      string    `json:"password"`
	token         string    `json:"token"`
	root          bool      `json:"root"`
	mail          string    `json:"mail"`
	displayname   string    `json:"displayname"`
	register_date time.Time `json:"register_date"`
	status        string    `json:"status"`
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
		err2 = resp.Scan(&user.username)
		if err2 != nil {
		}
		answers = append(answers, user.username)
	}
	if len(answers) == 1 {
		stmt, err = conn.Prepare("UPDATE inv_users SET token = ? WHERE username=?")
		if err != nil {
			panic(err)
		}
		token := utils.GenerateToken()
		stmt.Exec(token, answers[0])
		return true, token
	} else {
		return false, ""
	}

}
