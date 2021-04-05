package actions

import "github.com/MathisBurger/OpenInventory-Backend/database/models"

func AddRefreshToken(tkn *models.RefreshTokenModel) {
	conn := GetConn()
	defer conn.Close()
	stmt, _ := conn.Prepare("INSERT INTO `inv_refresh-token` (`ID`, `username`, `token`, `Deadline`) VALUES (NULL, ?, ?, ?);")
	defer stmt.Close()
	stmt.Exec(tkn.Username, tkn.Token, tkn.Deadline)
}