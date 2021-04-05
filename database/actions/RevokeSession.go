package actions

func RevokeSession(user string, token string) {
	conn := GetConn()
	defer conn.Close()
	stmt, _ := conn.Prepare("DELETE FROM `inv_refresh-token` WHERE `username`=? AND `token`=?")
	defer stmt.Close()
	stmt.Exec(user, token)
}