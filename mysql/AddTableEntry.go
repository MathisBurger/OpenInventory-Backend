package OwnSQL

import (
	"fmt"
	"github.com/MathisBurger/OpenInventory-Backend/config"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"strings"
)

type columnNameStruct struct {
	COLUMN_NAME string `json:"COLUMN_NAME"`
}

type Entries struct {
	Entries int `json:"entries"`
}

func AddTableEntry(displayname string, password string, token string, Tablename string, row map[string]interface{}) bool {
	perms := MySQL_loginWithToken(displayname, password, token)
	if !perms {
		return false
	} else {
		conn := GetConn()
		cfg, _ := config.ParseConfig()
		stmt, _ := conn.Prepare("select COLUMN_NAME from INFORMATION_SCHEMA.COLUMNS where TABLE_NAME=? and TABLE_SCHEMA=?;")
		resp, err := stmt.Query("table_"+Tablename, cfg.Db.Database)
		if err != nil {
			utils.LogError("[AddTableEntry.go, 29, SQL-StatementError] " + err.Error())
		}
		var columns []string
		for resp.Next() {
			var column columnNameStruct
			err = resp.Scan(&column.COLUMN_NAME)
			if err != nil {
				utils.LogError("[AddTableEntry.go, 36, SQL-ScanningError] " + err.Error())
			}
			if column.COLUMN_NAME != "id" {
				if row[column.COLUMN_NAME] != nil {
					columns = append(columns, column.COLUMN_NAME)
				} else {
					fmt.Println(column.COLUMN_NAME)
					defer resp.Close()
					defer stmt.Close()
					defer conn.Close()
					return false
				}
			}
		}
		stmt, _ = conn.Prepare("SELECT `min-perm-lvl` FROM `inv_tables` WHERE `name`=?;")
		type cacheStruct struct {
			MinPermLvl int `json:"min-perm-lvl"`
		}
		resp, err = stmt.Query(Tablename)
		if err != nil {
			utils.LogError("[AddTableEntry.go, 55, SQL-ScanningError] " + err.Error())
		}
		minPermLvl := 0
		for resp.Next() {
			var cache cacheStruct
			err = resp.Scan(&cache.MinPermLvl)
			if err != nil {
				utils.LogError("[AddTableEntry.go, 62, SQL-ScanningError] " + err.Error())
			}
			minPermLvl = cache.MinPermLvl
		}
		defer resp.Close()
		if CheckUserHasHigherPermission(conn, displayname, minPermLvl, "") {
			var builder strings.Builder
			builder.WriteString("INSERT INTO `table_" + Tablename + "`(`id`, ")
			for i, el := range columns {
				if i == (len(columns) - 1) {
					builder.WriteString("`" + el + "`")
				} else {
					builder.WriteString("`" + el + "`, ")
				}

			}
			builder.WriteString(") VALUES (NULL, ")
			for i, _ := range columns {
				if i == (len(columns) - 1) {
					builder.WriteString("?")
					break
				} else {
					builder.WriteString("?, ")
				}
			}
			builder.WriteString(");")
			stmt, err = conn.Prepare(builder.String())
			if err != nil {
				utils.LogError("[AddTableEntry.go, 73, SQL-StatementError] " + err.Error())
				defer conn.Close()
				return false
			}
			values := ParseToArray(row, columns)
			_, err = stmt.Exec(values...)
			if err != nil {
				utils.LogError("[AddTableEntry.go, 81, SQL-StatementError] " + err.Error())
				defer stmt.Close()
				defer conn.Close()
				return false
			}
			stmt, _ = conn.Prepare("SELECT `entries` FROM `inv_tables` WHERE `name`=?")
			resp, err = stmt.Query(Tablename)
			if err != nil {
				utils.LogError("[AddTableEntry.go, 89, SQL-StatementError] " + err.Error())
			}
			entries := 0
			for resp.Next() {
				var entry Entries
				err = resp.Scan(&entry.Entries)
				if err != nil {
					utils.LogError("[AddTableEntry.go, 96, SQL-ScanningError] " + err.Error())
				}
				entries = entry.Entries
			}
			entries += 1
			stmt, _ = conn.Prepare("UPDATE `inv_tables` SET `entries`=? WHERE `name`=?;")
			stmt.Exec(entries, Tablename)
			defer resp.Close()
			defer stmt.Close()
			defer conn.Close()
			return true
		} else {
			defer resp.Close()
			defer stmt.Close()
			defer conn.Close()
			return false
		}
	}
}
func ParseToArray(input map[string]interface{}, columns []string) []interface{} {
	v := make([]interface{}, len(input), len(input))
	idx := 0
	for _, value := range columns {
		v[idx] = input[value]
		idx++
	}
	return v
}
