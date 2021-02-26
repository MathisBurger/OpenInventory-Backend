package actions

import (
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"strings"
)

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
	exists, cols := SelectColumnScheme(Tablename)
	if !exists {
		return false
	}
	var columns []string
	for _, column := range cols {
		if column.COLUMN_NAME != "id" {
			if row[column.COLUMN_NAME] != nil {
				columns = append(columns, column.COLUMN_NAME)
			} else {
				return false
			}
		}
	}
	table := GetTableByName(Tablename)
	if CheckUserHasHigherPermission(conn, displayname, table.MinPermLvl, "") {
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
		stmt, err := conn.Prepare(builder.String())
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
		stmt, _ = conn.Prepare("UPDATE `inv_tables` SET `entries`=? WHERE `name`=?;")
		stmt.Exec(table.Entrys+1, Tablename)
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
