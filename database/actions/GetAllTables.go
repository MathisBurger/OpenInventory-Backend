package actions

import (
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
)

func GetAllTables(username string, password string, token string) []models.TableModel {
	status := MysqlLoginWithToken(username, password, token)
	if !status {
		return []models.TableModel{}
	}
	conn := GetConn()
	stmt, err := conn.Prepare("SELECT * FROM inv_tables")
	if err != nil {
		utils.LogError(err.Error(), "GetAllTables.go", 16)
	}
	resp, err := stmt.Query()
	if err != nil {
		utils.LogError(err.Error(), "GetAllTables.go", 20)
	}
	var tables []models.TableModel
	for resp.Next() {
		var table models.TableModel
		err = resp.Scan(&table.ID, &table.Name, &table.Entries, &table.MinPermLvl, &table.CreatedAt)
		if err != nil {
			panic(err)
		}
		if CheckUserHasHigherPermission(conn, username, table.MinPermLvl, "") {
			tables = append(tables, table)
		}
	}
	defer resp.Close()
	defer stmt.Close()
	defer conn.Close()
	return tables
}