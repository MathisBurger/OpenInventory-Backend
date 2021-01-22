package OwnSQL

import (
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
)

func GetAllTables(username string, password string, token string) []models.TableModel {
	status := MySQL_loginWithToken(username, password, token)
	if !status {
		return []models.TableModel{}
	} else {
		conn := GetConn()
		stmt, err := conn.Prepare("SELECT * FROM inv_tables")
		if err != nil {
			utils.LogError("[GetAllTables.go, 16, SQL-StatementError] " + err.Error())
		}
		resp, err := stmt.Query()
		if err != nil {
			utils.LogError("[GetAllTables.go, 20, SQL-StatementError] " + err.Error())
		}
		var tables []models.TableModel
		for resp.Next() {
			var table models.TableModel
			err = resp.Scan(&table.ID, &table.Name, &table.Entries, &table.CreatedAt)
			if err != nil {
				panic(err)
			}
			tables = append(tables, models.TableModel{table.ID, table.Name, table.Entries, table.CreatedAt})
		}
		defer resp.Close()
		defer stmt.Close()
		defer conn.Close()
		return tables
	}
}
