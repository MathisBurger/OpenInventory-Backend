package controller

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/models"
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
			utils.LogError(err.Error(), "AddUserController.go", 19)
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
	status := actions.MysqlLoginWithTokenRoot(obj.Username, obj.Password, obj.Token)
	if status {
		hash := utils.HashWithSalt(obj.User.Password)
		actions.AddUser(obj.User.Username, hash, obj.User.Root, obj.User.Mail, obj.User.Status)
		res, _ := models.GetJSONResponse("Successfully added user", "alert alert-success", "ok", "None", 200)
		return c.Send(res)
	}
	res, _ := models.GetJSONResponse("You do not have the permission to perform this", "alert alert-danger", "ok", "None", 200)
	return c.Send(res)
}

func checkAddUserRequest(obj models.AddUserRequestModel) bool {
	struct1 := models.AddUserStruct{"", "", false, "", ""}
	struct2 := models.AddUserStruct{"", "", true, "", ""}
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.User != struct1 && obj.User != struct2
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
	return strings.Contains(mail, "@") && len(strings.Split(mail, ".")) > 0
}
