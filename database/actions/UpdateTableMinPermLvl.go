package actions

////////////////////////////////////
// Updates perm level of table    //
////////////////////////////////////
func UpdateTableMinPermLvl(tablename string, newLvl int) {

	conn := GetConn()
	defer conn.Close()

	stmt, _ := conn.Prepare("UPDATE `inv_tables` SET `min-perm-lvl`=? WHERE `name`=?;")
	defer stmt.Close()

	stmt.Exec(newLvl, tablename)
}
