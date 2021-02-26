package actions

import (
	"github.com/MathisBurger/OpenInventory-Backend/utils"
)

type ColumnNameStruct2 struct {
	COLUMN_NAME string      `json:"COLUMN_NAME"`
	DATA_TYPE   string      `json:"DATA_TYPE"`
	MAX_LENGTH  interface{} `json:"CHARACTER_MAXIMUM_LENGTH"`
}

func GetTableColumns(displayname string, password string, token string, Tablename string) []Column {
	perms := MysqlLoginWithToken(displayname, password, token)
	if !perms {
		return []Column{}
	}
	conn := GetConn()
	stmt, err := conn.Prepare("SELECT `min-perm-lvl` FROM `inv_tables` WHERE `name`=?;")
	if err != nil {
		utils.LogError(err.Error(), "GetTableColumns.go", 22)
	}
	type cacheStruct struct {
		MinPermLvl int `json:"min-perm-lvl"`
	}
	resp, err := stmt.Query(Tablename)
	if err != nil {
		utils.LogError(err.Error(), "GetTableColumns.go", 29)
	}
	minPermLvl := 0
	for resp.Next() {
		var cache cacheStruct
		err = resp.Scan(&cache.MinPermLvl)
		if err != nil {
			utils.LogError(err.Error(), "GetTableColumns.go", 36)
		}
		minPermLvl = cache.MinPermLvl
	}
	if CheckUserHasHigherPermission(conn, displayname, minPermLvl, "") {
		exists, ans := SelectColumnScheme(Tablename)
		if !exists {
			return []Column{}
		}
		return ans
	}
	return []Column{}
}
