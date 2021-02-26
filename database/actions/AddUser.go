package actions

import "github.com/MathisBurger/OpenInventory-Backend/utils"

func AddUser(username string, hash string, root bool, mail string, status string) {
	conn := GetConn()
	stmt, err := conn.Prepare("INSERT INTO `inv_users` (`id`, `username`, `password`, `token`, `permissions`, `root`, `mail`, `displayname`, `register_date`, `status`) VALUES (NULL, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP(), ?);")
	if err != nil {
		utils.LogError(err.Error(), "AddUser.go", 9)
	}
	var perms string
	if root {
		perms = "default.everyone;default.root"
	} else {
		perms = "default.everyone"
	}
	stmt.Exec(username, hash, "None", perms, root, mail, username, status)
	defer stmt.Close()
	defer conn.Close()
}
