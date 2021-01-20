package OwnSQL

import (
	"fmt"
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
		stmt, _ := conn.Prepare("select COLUMN_NAME from INFORMATION_SCHEMA.COLUMNS where TABLE_NAME=?;")
		resp, err := stmt.Query("table_" + Tablename)
		if err != nil {
			panic(err.Error())
		}
		var columns []string
		for resp.Next() {
			var column columnNameStruct
			err = resp.Scan(&column.COLUMN_NAME)
			if err != nil {
				panic(err.Error())
			}
			columns = append(columns, column.COLUMN_NAME)
		}
		var builder strings.Builder
		builder.WriteString("INSERT INTO `table_" + Tablename + "`(`id`, ")
		for i, el := range columns {
			if el != "id" {
				if i == (len(columns) - 1) {
					builder.WriteString("`" + el + "`")
				} else {
					builder.WriteString("`" + el + "`, ")
				}
			}
		}
		builder.WriteString(") VALUES (NULL, ")
		fmt.Println("length:", len(columns))
		for i, _ := range columns {
			if i == (len(columns) - 2) {
				builder.WriteString("?")
				break
			} else {
				builder.WriteString("?, ")
			}
		}
		builder.WriteString(");")
		stmt, err = conn.Prepare(builder.String())
		if err != nil {
			defer conn.Close()
			return false
		}

		values := ParseToArray(row)
		_, err = stmt.Exec(values...)
		if err != nil {
			stmt.Close()
			conn.Close()
			return false
		}
		stmt, _ = conn.Prepare("SELECT `entries` FROM `inv_tables` WHERE `name`=?")
		resp, err = stmt.Query(Tablename)
		if err != nil {
			panic(err.Error())
		}
		entries := 0
		for resp.Next() {
			var entry Entries
			err = resp.Scan(&entry.Entries)
			if err != nil {
				panic(err.Error())
			}
			entries = entry.Entries
		}
		entries += 1
		stmt, _ = conn.Prepare("UPDATE `inv_tables` SET `entries`=? WHERE `name`=?;")
		stmt.Exec(entries, Tablename)
		resp.Close()
		stmt.Close()
		conn.Close()
		return true
	}
}
func ParseToArray(input map[string]interface{}) []interface{} {
	v := make([]interface{}, len(input), len(input))
	idx := 0
	for _, value := range input {
		v[idx] = value
		idx++
	}
	return v
}
