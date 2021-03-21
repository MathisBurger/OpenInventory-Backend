package actions

func UpdateUserRoor(root bool, username string) {
	conn := GetConn()
	defer conn.Close()
	stmt, _ := conn.Prepare("UPDATE `inv_users` SET `root`=? WHERE `username`=?")
	defer stmt.Close()
	stmt.Exec(root, username)
}
