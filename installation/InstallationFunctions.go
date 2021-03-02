package installation

import (
	"database/sql"
	"fmt"
	"github.com/MathisBurger/OpenInventory-Backend/config"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
)

// defined config path
var cfg = "./config/config.json"

// cache struct
type Table struct {
	name string `json:"name"`
}

// checks if config exists
func config_exists() bool {
	_, err := ioutil.ReadFile(cfg)
	return err == nil
}

// retuns config as string
func GetConfigContent() string {
	f, err := ioutil.ReadFile(cfg)
	if err != nil {
		return ""
	}
	return string(f)

}

// tests mysql connection
func TestMySQLConnection(cfg *config.Config) bool {

	connstr := cfg.Db.Username + ":" + cfg.Db.Password + "@tcp(" + cfg.Db.Host + ")/" + cfg.Db.Database

	conn, err := sql.Open("mysql", connstr)
	if err != nil {
		fmt.Println("Connection to database failed")
		return false
	}

	fmt.Println("Successfully connected to database")
	defer conn.Close()
	return true
}

// checks for all tables
func CheckForTables(cfg *config.Config) bool {

	// connect to database
	connstr := cfg.Db.Username + ":" + cfg.Db.Password + "@tcp(" + cfg.Db.Host + ")/" + cfg.Db.Database
	conn, err := sql.Open("mysql", connstr)
	if err != nil {
		return false
	}

	// statement preparing
	tables, err := conn.Query("SHOW TABLES LIKE 'inv_%';")
	if err != nil {
		return false
	}

	// fetching all tables
	var activeTables []string
	for tables.Next() {
		var table Table
		err = tables.Scan(&table.name)
		if err != nil {
			return false
		}
		activeTables = append(activeTables, table.name)
	}
	if len(activeTables) == 3 {
		fmt.Println("All required tables are existing")
		return true
	}

	// defined required tables
	requiredTables := [3]string{"inv_users", "inv_tables", "inv_permissions"}
	var outstandingTables []string

	// checking if table exists
	for _, el := range requiredTables {
		if !utils.ContainsStr(activeTables, el) {
			outstandingTables = append(outstandingTables, el)
			fmt.Println("Table", el, "does not exist")
		} else {
			fmt.Println("Table", el, "exists")
		}
	}

	// generate missing tables
	for _, tab := range outstandingTables {
		GenerateTable(conn, tab)
	}

	defer tables.Close()
	defer conn.Close()
	return true
}

// generates table
func GenerateTable(conn *sql.DB, name string) {
	// check table name
	switch name {
	case "inv_users":
		creationString := "CREATE TABLE inv_users(id INT(6) UNSIGNED AUTO_INCREMENT PRIMARY KEY, username VARCHAR(32), password VARCHAR(1024), token VARCHAR(32), permissions TEXT, root TINYINT(1), mail VARCHAR(128), displayname VARCHAR(32), register_date DATETIME, status VARCHAR(16));"
		conn.Exec(creationString)
		InsertDefaultUser(conn)
		fmt.Println("Created default user")
		break
	case "inv_tables":
		creationString := "CREATE TABLE inv_tables(`id` INT(6) UNSIGNED AUTO_INCREMENT PRIMARY KEY, `name` VARCHAR(32), `entries` INT(6), `min-perm-lvl` INT(6), `created_at` DATETIME);"
		conn.Exec(creationString)
		break
	case "inv_permissions":
		creationString := "CREATE TABLE `inv_permissions` ( `ID` INT NOT NULL AUTO_INCREMENT , `name` TEXT NOT NULL , `color` VARCHAR(11) NOT NULL , `permission-level` INT NOT NULL , PRIMARY KEY (`ID`))"
		conn.Exec(creationString)
		InsertDefaultPermissionGroups(conn)
		break

	}
	fmt.Println("created table", name)
}

// insert default user
func InsertDefaultUser(conn *sql.DB) {

	// "Admin123" as MD5
	hash := utils.HashWithSalt("e64b78fc3bc91bcbc7dc232ba8ec59e0")
	stmt, _ := conn.Prepare("INSERT INTO inv_users (id, username, password, token, permissions, root, mail, displayname, register_date, status) VALUES (NULL, 'root',  ?, 'None', 'default.everyone;default.root', '1', 'example@mail.de', 'root', current_timestamp(), 'enabled');")

	stmt.Exec(hash)
	defer stmt.Close()
}

// insert default permission groups
func InsertDefaultPermissionGroups(conn *sql.DB) {

	stmt, _ := conn.Prepare("INSERT INTO `inv_permissions` (`ID`, `name`, `color`, `permission-level`) VALUES (NULL, 'default.everyone', '96,97,98', '1');")
	defer stmt.Close()
	stmt.Exec()

	stmt, _ = conn.Prepare("INSERT INTO `inv_permissions` (`ID`, `name`, `color`, `permission-level`) VALUES (NULL, 'default.root', '96,97,98', '100');")
	defer stmt.Close()
	stmt.Exec()
}
