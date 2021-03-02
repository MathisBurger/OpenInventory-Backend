package actions

//////////////////////////////////////////
// Deletes an user from the database    //
//////////////////////////////////////////
func DeleteUser(username string) {
	conn := GetConn()
	defer conn.Close()

	stmt, _ := conn.Prepare("DELETE FROM `inv_users` WHERE `username`=?;")
	defer stmt.Close()

	stmt.Exec(username)
}
