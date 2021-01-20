package OwnSQL

import (
	"fmt"
)

type ColumnNameStruct2 struct {
	COLUMN_NAME string      `json:"COLUMN_NAME"`
	DATA_TYPE   string      `json:"DATA_TYPE"`
	MAX_LENGTH  interface{} `json:"CHARACTER_MAXIMUM_LENGTH"`
}

func GetTableColumns(displayname string, password string, token string, Tablename string) interface{} {
	perms := MySQL_loginWithToken(displayname, password, token)
	if !perms {
		return false
	} else {
		conn := GetConn()
		stmt, _ := conn.Prepare("select COLUMN_NAME, DATA_TYPE, CHARACTER_MAXIMUM_LENGTH from INFORMATION_SCHEMA.COLUMNS where TABLE_NAME=?;")
		resp, err := stmt.Query("table_" + Tablename)
		if err != nil {
			fmt.Println(err.Error())
			return false
		}
		var answers []ColumnNameStruct2
		for resp.Next() {
			var cache ColumnNameStruct2
			err = resp.Scan(&cache.COLUMN_NAME, &cache.DATA_TYPE, &cache.MAX_LENGTH)
			if err != nil {
				panic(err.Error())
			}
			answers = append(answers, cache)
		}
		defer resp.Close()
		defer stmt.Close()
		defer conn.Close()
		return answers
	}
}
