package actions

//////////////////////////////////////
// Changes number of table entrys   //
// by indicator                     //
//////////////////////////////////////
func ChangeNumOfEntrysBy(tablename string, indicator int) {

	table := GetTableByName(tablename)

	newNum := table.Entrys + indicator

	conn := GetConn()
	defer conn.Close()

	stmt, _ := conn.Prepare("UPDATE `inv_tables` SET `entries`=? WHERE `name`=?")
	defer stmt.Close()

	stmt.Exec(newNum, tablename)
}
