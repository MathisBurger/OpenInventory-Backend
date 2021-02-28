package actions

func RenameTableColumn(tablename string, newname string, datatype string, length string) bool {
	conn := GetConn()
	defer conn.Close()
	stmt, err := conn.Prepare("ALTER TABLE `table_" + tablename + "` CHANGE `" + tablename + "`  `" + newname + "` " + datatype +
		"(" + length + ") NULL DEFAULT NULL;")
	if err != nil {
		return false
	}
	_, err = stmt.Exec()
	return err != nil
}
