package actions

import (
	"database/sql"
	"github.com/MathisBurger/OpenInventory-Backend/config"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	_ "github.com/go-sql-driver/mysql"
)

func GetConn() (conn *sql.DB) {
	cfg, err := config.ParseConfig()
	if err != nil {
		utils.LogError(err.Error(), "utils.go", 13)
	}
	connstr := cfg.Db.Username + ":" + cfg.Db.Password + "@tcp(" + cfg.Db.Host + ")/" + cfg.Db.Database
	conn, err = sql.Open("mysql", connstr)
	if err != nil {
		utils.LogError(err.Error(), "utils.go", 18)
		return
	}
	return
}
