package actions

import "github.com/MathisBurger/OpenInventory-Backend/database/models"

func GetTableByName(name string) models.TableModel {
	conn := GetConn()
	defer conn.Close()
	stmt, _ := conn.Prepare("SELECT * FROM `inv_tables` WHERE `name`=?")
	defer stmt.Close()
	resp, _ := stmt.Query(name)
	defer resp.Close()
	resp.Next()
	return models.TableModel{}.Parse(resp)
}
