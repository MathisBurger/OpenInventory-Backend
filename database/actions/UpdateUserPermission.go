package actions

////////////////////////////////////
// Updates permission of user     //
////////////////////////////////////
func UpdateUserPermission(username string, permissions string) {

	conn := GetConn()
	defer conn.Close()

	stmt, _ := conn.Prepare("UPDATE `inv_users` SET `permissions`=? WHERE `username`=?;")
	defer stmt.Close()

	stmt.Exec(permissions, username)
}
