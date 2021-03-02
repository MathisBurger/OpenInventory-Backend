package actions

////////////////////////////////////
// Renames the given table        //
////////////////////////////////////
func RenameTable(oldname string, newname string) bool {

	conn := GetConn()
	defer conn.Close()

	stmt, _ := conn.Prepare("ALTER TABLE `table_" + oldname + "` RENAME `table_" + newname + "`;")
	defer stmt.Close()

	_, err := stmt.Exec()
	return err == nil
}
