package installation

import "fmt"
import "github.com/MathisBurger/OpenInventory-Backend/config"

func Install() bool {
	// check if config exists
	if !config_exists() {
		fmt.Println("Configuration file does not exist.  Please follow the instructions in the docs.")
		return false
	}
	fmt.Println("Configuration file exists")
	// get config content
	content := GetConfigContent()
	if content == "" {
		fmt.Println("Please fill your config")
		return false
	}
	fmt.Println("Config is not empty")
	cfg, err := config.ParseConfig()
	if err != nil {
		fmt.Println("Config syntax is wrong")
		return false
	}
	fmt.Println("Configuration parsed successfully")
	// test connection
	if TestMySQLConnection(cfg) {
		fmt.Println("Checking for tables...")
		// check for tables
		if CheckForTables(cfg) {
			return true
		}
		return false
	}
	return false

}
