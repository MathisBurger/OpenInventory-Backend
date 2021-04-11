package actions

func UpdateUser2FA(username string, status bool) {

	conn := GetConn()
	defer conn.Close()

	stmt, _ := conn.Prepare("UPDATE `inv_users` SET `2fa`=? WHERE `username`=?")
	defer stmt.Close()

	stmt.Exec(status, username)
}
