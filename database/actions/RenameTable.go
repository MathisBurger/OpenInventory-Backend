package actions

func RenameTable(oldname string, newname string) bool {
	conn := GetConn()
	defer conn.Close()
	stmt, err := conn.Prepare("ALTER TABLE `table_" + oldname + "` RENAME `table_" + newname + "`;")
	if err != nil {
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	return err == nil
}
