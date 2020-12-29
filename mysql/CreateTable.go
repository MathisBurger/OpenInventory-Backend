package OwnSQL

import (
	"fmt"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
)

func CreateTable(displayname string, password string, token string, Tablename string, RowConfig []models.RowConfigModel) bool {
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
				fmt.Println("hitler")
				fmt.Println(typeString)
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
			fmt.Println("Error with sql syntax")
			return false
		}
		stmt.Exec()
		stmt, _ = conn.Prepare("INSERT INTO `inv_tables` (`id`, `name`, `entries`, `created_at`) VALUES (NULL, ?, '0', current_timestamp);")
		stmt.Exec(Tablename)
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
		if utils.ContainsStr(names, el.Name) {
			return false
		}
	}
	return true
}
