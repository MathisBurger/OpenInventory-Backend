package actions

import (
	"database/sql"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"strings"
)

type cachePermissionLevelStruct struct {
	PermissionLevel int `json:"permission-level"`
}

func CheckUserHasHigherPermission(conn *sql.DB, username string, permLevel int, permName string) bool {
	highestPermission := GetHighestPermission(conn, username)
	if permLevel > 0 {
		return highestPermission >= permLevel
	} else if permName != "" {
		stmt, err := conn.Prepare("SELECT `permission-level` FROM `inv_permissions` WHERE `name`=?")
		defer stmt.Close()
		if err != nil {
			utils.LogError(err.Error(), "CheckHigherPermission.go", 21)
		}
		resp, err := stmt.Query(permName)
		defer resp.Close()
		if err != nil {
			utils.LogError(err.Error(), "CheckHigherPermission.go", 26)
		}
		wantedPermissionLevel := 0
		for resp.Next() {
			var cache cachePermissionLevelStruct
			err = resp.Scan(&cache.PermissionLevel)
			if err != nil {
				utils.LogError(err.Error(), "CheckHigherPermission.go", 33)
			}
			wantedPermissionLevel = cache.PermissionLevel
		}
		return highestPermission >= wantedPermissionLevel
	} else {
		return false
	}
}

func GetHighestPermission(conn *sql.DB, username string) int {
	_, user := GetUserByUsername(username)
	permissions := strings.Split(user.Permissions, ";")
	stmt, err := conn.Prepare("SELECT `permission-level` FROM `inv_permissions` WHERE `name`=?")
	if err != nil {
		utils.LogError(err.Error(), "CheckHigherPermission.go", 48)
	}
	highestPermission := 0
	for _, val := range permissions {
		resp, err := stmt.Query(val)
		if err != nil {
			utils.LogError(err.Error(), "CheckHigherPermission.go", 54)
		}
		for resp.Next() {
			var cache cachePermissionLevelStruct
			err = resp.Scan(&cache.PermissionLevel)
			if err != nil {
				utils.LogError(err.Error(), "CheckHigherPermission.go", 60)
			}
			if cache.PermissionLevel > highestPermission {
				highestPermission = cache.PermissionLevel
			}
		}
		defer resp.Close()
	}
	return highestPermission
}
