package actions

import "github.com/MathisBurger/OpenInventory-Backend/database/models"

func GetRefreshToken(token string) (bool, models.RefreshTokenModel) {
	
	conn := GetConn()
	defer conn.Close()
	stmt, _ := conn.Prepare("SELECT * FROM `inv_refresh-token` WHERE `token`=?")
	defer stmt.Close()
	resp, _ := stmt.Query(token)
	defer resp.Close()
	return models.RefreshTokenModel{}.Parse(resp)
}