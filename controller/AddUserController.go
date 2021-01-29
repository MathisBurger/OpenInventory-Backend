package controller

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	OwnSQL "github.com/MathisBurger/OpenInventory-Backend/mysql"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func AddUserController(c *fiber.Ctx) error {
	raw := string(c.Body())
	obj := models.AddUserRequestModel{}
	err := json.Unmarshal([]byte(raw), &obj)
	if err != nil {
		res, err := models.GetJSONResponse("Invaild JSON body", "alert alert-danger", "error", "None", 200)
		if err != nil {
			utils.LogError("[AddUserController.go, 19, InputError] " + err.Error())
		}
		return c.Send(res)
	}
	if !checkAddUserRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !checkUsernameLength(obj.User.Username) {
		res, _ := models.GetJSONResponse("This username is too long", "alert alert-warning", "ok", "None", 200)
		return c.Send(res)
	}
	if !checkPasswordLength(obj.User.Username) {
		res, _ := models.GetJSONResponse("This password is too long", "alert alert-warning", "ok", "None", 200)
		return c.Send(res)
	}
	if !checkEmail(obj.User.Mail) {
		res, _ := models.GetJSONResponse("Your mail is not a email-address", "alert alert-warning", "ok", "None", 200)
		return c.Send(res)
	}
	status := OwnSQL.MysqlLoginWithTokenRoot(obj.Username, obj.Password, obj.Token)
	hash := utils.HashWithSalt(obj.User.Password)
	if status {
		conn := OwnSQL.GetConn()
		stmt, err := conn.Prepare("INSERT INTO `inv_users` (`id`, `username`, `password`, `token`, `permissions`, `root`, `mail`, `displayname`, `register_date`, `status`) VALUES (NULL, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP(), ?);")
		if err != nil {
			utils.LogError("[AddUserController.gp, 45, SQL-StatementError] " + err.Error())
		}
		var perms string
		if obj.User.Root {
			perms = "default.everyone;default.root"
		} else {
			perms = "default.everyone"
		}
		stmt.Exec(obj.User.Username, hash, "None", perms, obj.User.Root, obj.User.Mail, obj.User.Username, obj.User.Status)
		defer stmt.Close()
		defer conn.Close()
		res, _ := models.GetJSONResponse("Successfully added user", "alert alert-success", "ok", "None", 200)
		return c.Send(res)
	} else {
		res, _ := models.GetJSONResponse("You do not have the permission to perform this", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
}

func checkAddUserRequest(obj models.AddUserRequestModel) bool {
	struct1 := models.AddUserStruct{"", "", false, "", ""}
	struct2 := models.AddUserStruct{"", "", true, "", ""}
	if obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.User != struct1 && obj.User != struct2 {
		return true
	} else {
		return false
	}
}

func checkUsernameLength(username string) bool {
	split := strings.Split(username, "")
	return len(split) < 32
}

func checkPasswordLength(hash string) bool {
	split := strings.Split(hash, "")
	return len(split) < 1024
}

func checkEmail(mail string) bool {
	if strings.Contains(mail, "@") && len(strings.Split(mail, ".")) > 0 {
		return true
	} else {
		return false
	}
}
