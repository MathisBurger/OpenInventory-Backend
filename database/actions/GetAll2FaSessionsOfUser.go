package actions

import "github.com/MathisBurger/OpenInventory-Backend/database/models"

func GetAll2FaSessionsOfUser(owner string) []models.TwoFactorModel {

	conn := GetConn()
	defer conn.Close()

	stmt, _ := conn.Prepare("SELECT * FROM `inv_2fa-sessions` WHERE `owner`=?")
	defer stmt.Close()

	resp, _ := stmt.Query(owner)
	defer resp.Close()

	return models.TwoFactorModel{}.ParseAll(resp)
}
