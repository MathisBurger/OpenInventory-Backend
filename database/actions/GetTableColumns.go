package actions

import "github.com/MathisBurger/OpenInventory-Backend/database/models"

/////////////////////////////////////////////
// Queries all columns of specific table   //
/////////////////////////////////////////////
func GetTableColumns(displayname string, password string, token string, Tablename string) []Column {

	perms := MysqlLoginWithToken(displayname, password, token)

	if !perms {
		return []Column{}
	}

	conn := GetConn()
	stmt, _ := conn.Prepare("SELECT `min-perm-lvl` FROM `inv_tables` WHERE `name`=?;")

	// defining struct for later use
	type cacheStruct struct {
		MinPermLvl int `json:"min-perm-lvl"`
	}

	resp, _ := stmt.Query(Tablename)

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
