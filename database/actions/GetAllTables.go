package actions

import (
	"github.com/MathisBurger/OpenInventory-Backend/models"
)

////////////////////////////////////
// Queries all tables             //
////////////////////////////////////
func GetAllTables(username string, password string, token string) []models.TableModel {
	status := MysqlLoginWithToken(username, password, token)

	if !status {
		return []models.TableModel{}
	}

	conn := GetConn()
	defer conn.Close()

	stmt, _ := conn.Prepare("SELECT * FROM inv_tables")
	defer stmt.Close()

	resp, _ := stmt.Query()
	defer resp.Close()

	var tables []models.TableModel

	// fetch values into array
	for resp.Next() {
		var table models.TableModel

		_ = resp.Scan(&table.ID, &table.Name, &table.Entries, &table.MinPermLvl, &table.CreatedAt)

		if CheckUserHasHigherPermission(conn, username, table.MinPermLvl, "") {
			tables = append(tables, table)
		}
	}

	return tables
}
