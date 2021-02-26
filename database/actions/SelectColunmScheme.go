package actions

import (
	"github.com/MathisBurger/OpenInventory-Backend/config"
)

type Column struct {
	COLUMN_NAME string      `json:"COLUMN_NAME"`
	DATA_TYPE   string      `json:"DATA_TYPE"`
	MAX_LENGTH  interface{} `json:"CHARACTER_MAXIMUM_LENGTH"`
}

func SelectColumnScheme(tablename string) (bool, []Column) {
	conn := GetConn()
	stmt, _ := conn.Prepare("select COLUMN_NAME, DATA_TYPE, CHARACTER_MAXIMUM_LENGTH from INFORMATION_SCHEMA.COLUMNS where TABLE_NAME=? and TABLE_SCHEMA=?;")
	defer stmt.Close()
	cfg, _ := config.ParseConfig()
	resp, err := stmt.Query("table_"+tablename, cfg.Db.Database)
	defer resp.Close()
	if err != nil {
		return false, nil
	}
	var answers []Column
	for resp.Next() {
		var cache Column
		_ = resp.Scan(&cache.COLUMN_NAME, &cache.DATA_TYPE, &cache.MAX_LENGTH)
		answers = append(answers, cache)
	}
	return true, answers
}
