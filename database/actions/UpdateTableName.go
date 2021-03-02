package actions

///////////////////////////////////
// Updates tablename in internal //
// inv_tables table              //
///////////////////////////////////
func UpdateTablename(old string, new string) {

	conn := GetConn()
	defer conn.Close()

	stmt, _ := conn.Prepare("UPDATE `inv_tables` SET `name`=? WHERE `name`=?")
	defer stmt.Close()

	stmt.Exec(new, old)
}
