package actions

import (
	"database/sql"
	"github.com/MathisBurger/OpenInventory-Backend/database/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"strings"
)

///////////////////////////////////////////////
// Checks if user has higher permission      //
// than other user                           //
///////////////////////////////////////////////
func CheckUserHasHigherPermission(conn *sql.DB, username string, permLevel int, permName string) bool {
	// gets value from function below
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

		wantedPermissionLevel := models.PermissionModel{}.ParseAll(resp)[0].PermissionLevel

		return highestPermission >= wantedPermissionLevel
	} else {
		return false
	}
}

///////////////////////////////////////////
// returns highest permission of         //
// given user                            //
///////////////////////////////////////////
func GetHighestPermission(conn *sql.DB, username string) int {

	_, user := GetUserByUsername(username)

	// all permission groups of user
	permissions := strings.Split(user.Permissions, ";")

	// get highest permission of user
	highestPermission := 0
	for _, val := range permissions {

		_, perm := GetPermissionByName(val)

		if perm.PermissionLevel > highestPermission {
			highestPermission = perm.PermissionLevel
		}
	}
	return highestPermission
}
