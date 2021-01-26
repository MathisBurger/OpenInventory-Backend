package OwnSQL

import (
	"github.com/MathisBurger/OpenInventory-Backend/config"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
)

type ColumnNameStruct2 struct {
	COLUMN_NAME string      `json:"COLUMN_NAME"`
	DATA_TYPE   string      `json:"DATA_TYPE"`
	MAX_LENGTH  interface{} `json:"CHARACTER_MAXIMUM_LENGTH"`
}

func GetTableColumns(displayname string, password string, token string, Tablename string) []ColumnNameStruct2 {
	perms := MySQL_loginWithToken(displayname, password, token)
	if !perms {
		return []ColumnNameStruct2{}
	} else {
		conn := GetConn()
		stmt, err := conn.Prepare("SELECT `min-perm-lvl` FROM `inv_tables` WHERE `name`=?;")
		if err != nil {
			utils.LogError("[GetTableColumns.go, 21, SQL-StatementError] " + err.Error())
		}
		type cacheStruct struct {
			MinPermLvl int `json:"min-perm-lvl"`
		}
		resp, err := stmt.Query(Tablename)
		if err != nil {
			utils.LogError("[GetTableColumns.go, 28, SQL-StatementError] " + err.Error())
		}
		minPermLvl := 0
		for resp.Next() {
			var cache cacheStruct
			err = resp.Scan(&cache.MinPermLvl)
			if err != nil {
				utils.LogError("[GetTableColumns.go, 35, SQL-StatementError] " + err.Error())
			}
			minPermLvl = cache.MinPermLvl
		}
		if CheckUserHasHigherPermission(conn, displayname, minPermLvl, "") {
			cfg, _ := config.ParseConfig()
			stmt, _ = conn.Prepare("select COLUMN_NAME, DATA_TYPE, CHARACTER_MAXIMUM_LENGTH from INFORMATION_SCHEMA.COLUMNS where TABLE_NAME=? and TABLE_SCHEMA=?;")
			resp, err := stmt.Query("table_"+Tablename, cfg.Db.Database)
			if err != nil {
				utils.LogError("[GetTableColumns.go, 43, SQL-StatementError] " + err.Error())
				return []ColumnNameStruct2{}
			}
			var answers []ColumnNameStruct2
			for resp.Next() {
				var cache ColumnNameStruct2
				err = resp.Scan(&cache.COLUMN_NAME, &cache.DATA_TYPE, &cache.MAX_LENGTH)
				if err != nil {
					utils.LogError("[GetTableColumns.go, 51, SQL-ScanningError] " + err.Error())
				}
				answers = append(answers, cache)
			}
			defer resp.Close()
			defer stmt.Close()
			defer conn.Close()
			return answers
		} else {
			return []ColumnNameStruct2{}
		}
	}
}
