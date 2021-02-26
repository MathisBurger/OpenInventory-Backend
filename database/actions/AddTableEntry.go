package actions

import (
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
	perms := MysqlLoginWithToken(displayname, password, token)
	if !perms {
		return false
	}
	conn := GetConn()
	defer conn.Close()
	cfg, _ := config.ParseConfig()
	stmt, _ := conn.Prepare("select COLUMN_NAME from INFORMATION_SCHEMA.COLUMNS where TABLE_NAME=? and TABLE_SCHEMA=?;")
	defer stmt.Close()
	resp, err := stmt.Query("table_"+Tablename, cfg.Db.Database)
	defer resp.Close()
	if err != nil {
		utils.LogError(err.Error(), "AddTableEntry.go", 31)
	}
	var columns []string
	for resp.Next() {
		var column columnNameStruct
		err = resp.Scan(&column.COLUMN_NAME)
		if err != nil {
			utils.LogError(err.Error(), "AddTableEntry.go", 38)
		}
		if column.COLUMN_NAME != "id" {
			if row[column.COLUMN_NAME] != nil {
				columns = append(columns, column.COLUMN_NAME)
			} else {
				return false
			}
		}
	}
	stmt, _ = conn.Prepare("SELECT `min-perm-lvl` FROM `inv_tables` WHERE `name`=?;")
	defer stmt.Close()
	type cacheStruct struct {
		MinPermLvl int `json:"min-perm-lvl"`
	}
	resp, err = stmt.Query(Tablename)
	defer resp.Close()
	if err != nil {
		utils.LogError(err.Error(), "AddTableEntry.go", 55)
	}
	minPermLvl := 0
	for resp.Next() {
		var cache cacheStruct
		err = resp.Scan(&cache.MinPermLvl)
		if err != nil {
			utils.LogError(err.Error(), "AddTableEntry.go", 63)
		}
		minPermLvl = cache.MinPermLvl
	}
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
		for i := range columns {
			if i == (len(columns) - 1) {
				builder.WriteString("?")
				break
			} else {
				builder.WriteString("?, ")
			}
		}
		builder.WriteString(");")
		stmt, err = conn.Prepare(builder.String())
		defer stmt.Close()
		if err != nil {
			utils.LogError(err.Error(), "AddTableEntry.go", 91)
			return false
		}
		values := ParseToArray(row, columns)
		_, err = stmt.Exec(values...)
		if err != nil {
			utils.LogError(err.Error(), "AddTableEntry.go", 97)
			return false
		}
		stmt, _ = conn.Prepare("SELECT `entries` FROM `inv_tables` WHERE `name`=?")
		defer stmt.Close()
		resp, err = stmt.Query(Tablename)
		defer resp.Close()
		if err != nil {
			utils.LogError(err.Error(), "AddTableEntry.go", 105)
		}
		entries := 0
		for resp.Next() {
			var entry Entries
			err = resp.Scan(&entry.Entries)
			if err != nil {
				utils.LogError(err.Error(), "AddTableEntry.go", 112)
			}
			entries = entry.Entries
		}
		entries++
		stmt, _ = conn.Prepare("UPDATE `inv_tables` SET `entries`=? WHERE `name`=?;")
		stmt.Exec(entries, Tablename)
		defer stmt.Close()
		return true
	}
	return false
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
