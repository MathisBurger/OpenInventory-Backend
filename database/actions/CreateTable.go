package actions

import (
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"strings"
)

////////////////////////////////////
// Creates a table                //
////////////////////////////////////
func CreateTable(displayname string, password string, token string, Tablename string, RowConfig []models.RowConfigModel, MinPermLvl int) bool {

	// check user login
	if !MysqlLoginWithToken(displayname, password, token) {
		return false
	}

	if !CheckColumnNames(RowConfig) {
		return false
	}

	cache := ""

	// check if columns are valid
	// add build sql string
	for _, row := range RowConfig {
		typeString := checkType(row)
		if typeString == "" {
			return false
		}
		if strings.Compare(row.Name, "") == 0 {
			return false
		}
		cache += typeString
	}

	// removing last char from array
	chars := []rune(cache)
	index := len(chars) - 1

	finStr := ""
	// building final sql string from modified char array
	for i, el := range chars {
		if i == index {
			break
		} else {
			finStr += string(el)
		}
	}

	creationString := "CREATE TABLE IF NOT EXISTS `table_" + Tablename + "` (id INT(6) UNSIGNED AUTO_INCREMENT PRIMARY KEY, " + finStr + ");"

	conn := GetConn()
	defer conn.Close()

	if !CheckUserHasHigherPermission(conn, displayname, MinPermLvl, "") {
		return false
	}
	stmt, err := conn.Prepare(creationString)
	defer stmt.Close()
	if err != nil {
		utils.LogError(err.Error(), "CreateTable.go", 47)
		return false
	}

	stmt.Exec()

	stmt, _ = conn.Prepare("INSERT INTO `inv_tables` (`id`, `name`, `entries`, `min-perm-lvl`, `created_at`) VALUES (NULL, ?, '0', ?, current_timestamp);")
	defer stmt.Close()

	stmt.Exec(Tablename, MinPermLvl)

	return true
}

///////////////////////////////////
// returns sql specific string   //
// based on rowType              //
///////////////////////////////////
func checkType(row models.RowConfigModel) string {
	switch row.Type {
	case "INT":
		return "`" + row.Name + "` INT(11) NOT NULL,"
	case "FLOAT":
		return "`" + row.Name + "` float NOT NULL,"
	case "BOOLEAN":
		return "`" + row.Name + "` tinyint(1) NOT NULL,"
	case "TEXT":
		return "`" + row.Name + "` TEXT NOT NULL,"
	default:
		return ""
	}
}

///////////////////////////////////////
// checks if column names are valid  //
///////////////////////////////////////
func CheckColumnNames(columns []models.RowConfigModel) bool {
	var names []string
	for _, el := range columns {
		if utils.ContainsStr(names, el.Name) || utils.ContainsStr(names, "id") || utils.ContainsStr(names, "ID") {
			return false
		}

	}
	return true
}
