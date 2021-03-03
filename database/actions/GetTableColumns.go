package actions

import (
	"github.com/MathisBurger/OpenInventory-Backend/database/models"
)

/////////////////////////////////////////////
// Queries all columns of specific table   //
/////////////////////////////////////////////
func GetTableColumns(displayname string, password string, token string, Tablename string) []Column {

	perms := MysqlLoginWithToken(displayname, password, token)

	if !perms {
		return []Column{}
	}

	conn := GetConn()
	defer conn.Close()

	stmt, _ := conn.Prepare("SELECT * FROM `inv_tables` WHERE `name`=?;")
	defer stmt.Close()

	resp, _ := stmt.Query(Tablename)
	defer resp.Close()

	minPermLvl := models.TableModel{}.ParseAll(resp)[0].MinPermLvl

	// check for higher permission
	if CheckUserHasHigherPermission(conn, displayname, minPermLvl, "") {
		exists, ans := SelectColumnScheme(Tablename)
		if !exists {
			return []Column{}
		}
		return ans
	}

	return []Column{}
}
