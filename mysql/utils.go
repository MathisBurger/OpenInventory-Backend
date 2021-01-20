package OwnSQL

import (
	"database/sql"
	"github.com/MathisBurger/OpenInventory-Backend/config"
	_ "github.com/go-sql-driver/mysql"
)

func GetConn() (conn *sql.DB) {
	cfg, err := config.ParseConfig()
	if err != nil {
		panic(err)
	}
	connstr := cfg.Db.Username + ":" + cfg.Db.Password + "@tcp(" + cfg.Db.Host + ")/" + cfg.Db.Database
	conn, err = sql.Open("mysql", connstr)
	if err != nil {
		panic(err.Error())
		return
	} else {
		return
	}
}
