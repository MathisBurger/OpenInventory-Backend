package OwnSQL

import (
	"database/sql"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"strings"
)

func CheckUserHasHigherPermission(conn *sql.DB, username string, permLevel int, permName string) bool {
	stmt, err := conn.Prepare("SELECT permissions FROM `inv_users` WHERE `username`=?;")
	if err != nil {
		utils.LogError("[CheckHigherPermission.go, 11, SQL-StatementError] " + err.Error())
	}
	resp, err := stmt.Query(username)
	if err != nil {
		utils.LogError("[CheckHigherPermission.go, 15, SQL-StatementError] " + err.Error())
	}
	type cacheStruct struct {
		Permissions string `json:"permissions"`
	}
	var permissions []string
	for resp.Next() {
		var cache cacheStruct
		err = resp.Scan(&cache.Permissions)
		if err != nil {
			utils.LogError("[CheckHigherPermission.go, 25, SQL-StatementError] " + err.Error())
		}
		permissions = strings.Split(cache.Permissions, ";")
	}
	defer resp.Close()
	stmt, err = conn.Prepare("SELECT `permission-level` FROM `inv_permissions` WHERE `name`=?")
	if err != nil {
		utils.LogError("[CheckHigherPermission.go, 33, SQL-StatementError] " + err.Error())
	}
	highestPermission := 0
	type cachePermissionLevelStruct struct {
		PermissionLevel int `json:"permission-level"`
	}
	for _, val := range permissions {
		resp, err = stmt.Query(val)
		if err != nil {
			utils.LogError("[CheckHigherPermission.go, 39, SQL-StatementError] " + err.Error())
		}
		for resp.Next() {
			var cache cachePermissionLevelStruct
			err = resp.Scan(&cache.PermissionLevel)
			if err != nil {
				utils.LogError("[CheckHigherPermission.go, 48, SQL-StatementError] " + err.Error())
			}
			if cache.PermissionLevel > highestPermission {
				highestPermission = cache.PermissionLevel
			}
		}
		defer resp.Close()
	}
	if permLevel > 0 {
		defer stmt.Close()
		return highestPermission >= permLevel
	} else if permName != "" {
		stmt, err = conn.Prepare("SELECT `permission-level` FROM `inv_permissions` WHERE `name`=?")
		if err != nil {
			utils.LogError("[CheckHigherPermission.go, 61, SQL-StatementError] " + err.Error())
		}
		resp, err = stmt.Query(permName)
		if err != nil {
			utils.LogError("[CheckHigherPermission.go, 65, SQL-StatementError] " + err.Error())
		}
		wantedPermissionLevel := 0
		for resp.Next() {
			var cache cachePermissionLevelStruct
			err = resp.Scan(&cache.PermissionLevel)
			if err != nil {
				utils.LogError("[CheckHigherPermission.go, 72, SQL-StatementError] " + err.Error())
			}
			wantedPermissionLevel = cache.PermissionLevel
		}
		defer resp.Close()
		defer stmt.Close()
		return highestPermission >= wantedPermissionLevel
	} else {
		utils.LogError("[CheckHigherPermission.go, 81, InputError] " + err.Error())
		return false
	}
}
