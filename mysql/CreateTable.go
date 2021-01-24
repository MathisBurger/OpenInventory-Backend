package OwnSQL

import (
	"fmt"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"strings"
)

func CreateTable(displayname string, password string, token string, Tablename string, RowConfig []models.RowConfigModel, MinPermLvl int) bool {
	perms := MySQL_loginWithToken(displayname, password, token)
	if !perms {
		return false
	} else {
		if !CheckColumnNames(RowConfig) {
			return false
		}
		cache := ""
		for _, row := range RowConfig {
			typeString := checkType(row)
			if typeString == "" {
				fmt.Println(typeString)
				return false
			}
			if strings.Compare(row.Name, "") == 0 {
				return false
			}
			cache += typeString
		}
		chars := []rune(cache)
		index := len(chars) - 1
		finStr := ""
		for i, el := range chars {
			if i == index {
				break
			} else {
				finStr += string(el)
			}
		}
		creationString := "CREATE TABLE IF NOT EXISTS `table_" + Tablename + "` (id INT(6) UNSIGNED AUTO_INCREMENT PRIMARY KEY, " + finStr + ");"
		conn := GetConn()
		stmt, err := conn.Prepare(creationString)
		if err != nil {
			utils.LogError("[CreateTable.go, 44, SQL-StatementError] " + err.Error())
			return false
		}
		stmt.Exec()
		stmt, _ = conn.Prepare("INSERT INTO `inv_tables` (`id`, `name`, `entries`, `min-perm-lvl`, `created_at`) VALUES (NULL, ?, '0', ?, current_timestamp);")
		stmt.Exec(Tablename, MinPermLvl)
		defer stmt.Close()
		defer conn.Close()
		return true
	}
}

func checkType(row models.RowConfigModel) string {
	switch row.Type {
	case "INT":
		return "`" + row.Name + "` INT(11),"
	case "FLOAT":
		return "`" + row.Name + "` float,"
	case "BOOLEAN":
		return "`" + row.Name + "` tinyint(1),"
	case "STRING8Chars":
		return "`" + row.Name + "` VARCHAR(8),"
	case "STRING16Chars":
		return "`" + row.Name + "` VARCHAR(16),"
	case "STRING64Chars":
		return "`" + row.Name + "` VARCHAR(64),"
	case "STRING128Chars":
		return "`" + row.Name + "` VARCHAR(128),"
	case "STRING1024Chars":
		return "`" + row.Name + "` VARCHAR(1024),"
	default:
		fmt.Println("incorrect type:", row.Type)
		return ""
	}
}

func CheckColumnNames(columns []models.RowConfigModel) bool {
	var names []string
	for _, el := range columns {
		if utils.ContainsStr(names, el.Name) || utils.ContainsStr(names, "id") || utils.ContainsStr(names, "ID") {
			return false
		}

	}
	return true
}
