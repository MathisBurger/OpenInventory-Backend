package actions

import (
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"strings"
)

type Entries struct {
	Entries int `json:"entries"`
}

/////////////////////////////////////
// Adds entry to given table       //
// returns status of action        //
/////////////////////////////////////
func AddTableEntry(displayname string, password string, token string, Tablename string, row map[string]interface{}) bool {

	// check login
	if !MysqlLoginWithToken(displayname, password, token) {
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

		// do not add id as parameter
		if column.COLUMN_NAME != "id" {

			// row should contain column
			if row[column.COLUMN_NAME] != nil {
				columns = append(columns, column.COLUMN_NAME)
			} else {
				return false
			}
		}
	}
	table := GetTableByName(Tablename)

	// check permission
	if CheckUserHasHigherPermission(conn, displayname, table.MinPermLvl, "") {

		// create insertion string
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

		// prepare statement from string builder
		stmt, err := conn.Prepare(builder.String())
		defer stmt.Close()

		if err != nil {
			utils.LogError(err.Error(), "AddTableEntry.go", 91)
			return false
		}

		// get values array from row
		values := ParseToArray(row, columns)

		_, err = stmt.Exec(values...)

		if err != nil {
			utils.LogError(err.Error(), "AddTableEntry.go", 97)
			return false
		}

		ChangeNumOfEntrysBy(Tablename, 1)

		return true
	}
	return false
}

// parse map to value array
// needs static map and column names
func ParseToArray(input map[string]interface{}, columns []string) []interface{} {
	v := make([]interface{}, len(input), len(input))
	idx := 0
	for _, value := range columns {
		v[idx] = input[value]
		idx++
	}
	return v
}
