package controller

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	OwnSQL "github.com/MathisBurger/OpenInventory-Backend/mysql"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

func AddUserController(c *fiber.Ctx) error {
	raw := string(c.Body())
	obj := models.AddUserRequestModel{}
	err := json.Unmarshal([]byte(raw), &obj)
	if err != nil {
		response, err := models.GetJsonResponse("Invaild JSON body", "alert alert-danger", "error", "None", 200)
		if err != nil {
			panic(err)
		}
		return c.Send(response)
	}
	if !checkAddUserRequest(obj) {
		resp, _ := models.GetJsonResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(resp)
	}
	status := OwnSQL.MySQL_loginWithToken_ROOT(obj.Username, obj.Password, obj.Token)
	if status {
		conn := OwnSQL.GetConn()
		stmt, err := conn.Prepare("INSERT INTO `inv_users` (`id`, `username`, `password`, `token`, `root`, `mail`, `displayname`, `register_date`, `status`) VALUES (NULL, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP(), ?);")
		if err != nil {
			panic(err.Error())
		}
		stmt.Exec(obj.User.Username, utils.HashWithSalt(obj.User.Password), "None", obj.User.Root, obj.User.Mail, obj.User.Username, obj.User.Status)
		stmt.Close()
		conn.Close()
		resp, _ := models.GetJsonResponse("Successfully added user", "alert alert-success", "ok", "None", 200)
		return c.Send(resp)
	} else {
		resp, _ := models.GetJsonResponse("You do not have the permission to perform this", "alert alert-danger", "ok", "None", 200)
		return c.Send(resp)
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
