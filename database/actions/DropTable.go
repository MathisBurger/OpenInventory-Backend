package actions

func DropTable(tablename string) {
	conn := GetConn()
	defer conn.Close()
	stmt, _ := conn.Prepare("DROP TABLE `table_" + tablename + "`;")
	defer stmt.Close()
	stmt.Exec()
	stmt, _ = conn.Prepare("DELETE FROM `inv_tables` WHERE `name`=?")
	defer stmt.Close()
	stmt.Exec(tablename)
}
