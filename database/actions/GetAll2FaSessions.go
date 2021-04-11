package actions

import "github.com/MathisBurger/OpenInventory-Backend/database/models"

func GetAll2FaSessions() []models.TwoFactorModel {

	conn := GetConn()
	defer conn.Close()

	stmt, _ := conn.Prepare("SELECT * FROM `inv_2fa-sessions`")
	defer stmt.Close()

	resp, _ := stmt.Query()
	defer resp.Close()

	return models.TwoFactorModel{}.ParseAll(resp)
}
