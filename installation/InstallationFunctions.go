package installation

import (
	"database/sql"
	"fmt"
	"github.com/MathisBurger/OpenInventory-Backend/config"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
)

var cfg = "./config/config.json"

type Table struct {
	name string `json:"name"`
}

func config_exists() bool {
	_, err := ioutil.ReadFile(cfg)
	if err != nil {
		return false
	} else {
		return true
	}
}

func GetConfigContent() string {
	f, err := ioutil.ReadFile(cfg)
	if err != nil {
		return ""
	} else {
		return string(f)
	}
}

func TestMySQLConnection(cfg *config.Config) bool {
	connstr := cfg.Db.Username + ":" + cfg.Db.Password + "@tcp(" + cfg.Db.Host + ")/" + cfg.Db.Database
	conn, err := sql.Open("mysql", connstr)
	if err != nil {
		fmt.Println("Connection to database failed")
		defer conn.Close()
		return false
	} else {
		fmt.Println("Successfully connected to database")
		defer conn.Close()
		return true
	}
}

func CheckForTables(cfg *config.Config) bool {
	connstr := cfg.Db.Username + ":" + cfg.Db.Password + "@tcp(" + cfg.Db.Host + ")/" + cfg.Db.Database
	conn, err := sql.Open("mysql", connstr)
	if err != nil {
		return false
	}
	tables, err := conn.Query("SHOW TABLES LIKE 'inv_%';")
	if err != nil {
		return false
	}
	var activeTables []string
	for tables.Next() {
		var table Table
		err = tables.Scan(&table.name)
		if err != nil {
			return false
		}
		activeTables = append(activeTables, table.name)
	}
	if len(activeTables) == 2 {
		fmt.Println("All required tables are existing")
		return true
	}
	requiredTables := [2]string{"inv_users", "inv_tables"}
	var outstandingTables []string
	for _, el := range requiredTables {
		if !utils.ContainsStr(activeTables, el) {
			outstandingTables = append(outstandingTables, el)
			fmt.Println("Table", el, "does not exist")
		} else {
			fmt.Println("Table", el, "exists")
		}
	}
	for _, tab := range outstandingTables {
		GenerateTable(conn, tab)
	}
	tables.Close()
	defer conn.Close()
	return true
}

func GenerateTable(conn *sql.DB, name string) {
	switch name {
	case "inv_users":
		creationString := "CREATE TABLE inv_users(id INT(6) UNSIGNED AUTO_INCREMENT PRIMARY KEY, username VARCHAR(32), password VARCHAR(1024), token VARCHAR(32), root TINYINT(1), mail VARCHAR(128), displayname VARCHAR(32), register_date DATETIME, status VARCHAR(16));"
		conn.Exec(creationString)
		InsertDefaultUser(conn)
		fmt.Println("Created default user")
		break
	case "inv_tables":
		creationString := "CREATE TABLE inv_tables(id INT(6) UNSIGNED AUTO_INCREMENT PRIMARY KEY, name VARCHAR(32), entries INT(6), created_at DATETIME);"
		conn.Exec(creationString)
		break
	}
	fmt.Println("created table", name)
}

func InsertDefaultUser(conn *sql.DB) {
	hash := utils.HashWithSalt("Admin123")
	stmt, err := conn.Prepare("INSERT INTO inv_users (id, username, password, token, root, mail, displayname, register_date, status) VALUES (NULL, 'root', ?, 'None', '1', 'example@mail.de', 'root', current_timestamp(), 'enabled');")
	if err != nil {
		panic(err.Error())
	}
	stmt.Exec(hash)
	defer stmt.Close()
}
