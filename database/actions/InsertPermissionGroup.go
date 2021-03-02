package actions

////////////////////////////////
// Creates permission group   //
////////////////////////////////
func InsertPermissionGroup(name string, color string, permLevel int) {
	conn := GetConn()
	defer conn.Close()
	stmt, _ := conn.Prepare("INSERT INTO `inv_permissions` (`ID`, `name`, `color`, `permission-level`) VALUES (NULL, ?, ?, ?);")
	defer stmt.Close()
	stmt.Exec(name, color, permLevel)
}
