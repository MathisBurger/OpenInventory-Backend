package actions

func DeleteTableEntry(entryID int, tablename string) bool {
	conn := GetConn()
	defer conn.Close()
	stmt, _ := conn.Prepare("DELETE FROM `table_" + tablename + "` WHERE `id`=?")
	defer stmt.Close()
	aff, _ := stmt.Exec(entryID)
	rowsAffected, _ := aff.RowsAffected()
	return rowsAffected != 0
}
