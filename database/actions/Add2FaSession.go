package actions

func Add2FaSession(secret string, owner string) {

	conn := GetConn()
	defer conn.Close()

	stmt, _ := conn.Prepare("INSERT INTO `inv_2fa-sessions` (`ID`, `secret`, `owner`) VALUES (NULL, ?, ?);")
	defer stmt.Close()

	stmt.Exec(secret, owner)
}
