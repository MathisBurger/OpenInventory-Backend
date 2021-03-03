package actions

////////////////////////////////////
// Renames a column of table      //
////////////////////////////////////
func RenameTableColumn(tablename string, oldname string, newname string, datatype string, length string) bool {

	conn := GetConn()
	defer conn.Close()

	stmt, err := conn.Prepare("ALTER TABLE `table_" + tablename + "` CHANGE `" + oldname + "`  `" + newname + "` " + datatype +
		"(" + length + ") NULL DEFAULT NULL;")
	if err != nil {
		return false
	}

	_, err = stmt.Exec()
	return err == nil
}
