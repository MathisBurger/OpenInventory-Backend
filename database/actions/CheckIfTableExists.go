package actions

//////////////////////////////////////////
// Checks if table exists               //
//////////////////////////////////////////
func CheckIfTableExists(name string) bool {
	conn := GetConn()
	defer conn.Close()

	stmt, _ := conn.Prepare("SELECT * FROM `inv_tables` WHERE `name`=?")
	defer stmt.Close()

	resp, _ := stmt.Query(name)
	defer resp.Close()

	return resp.Next()
}
