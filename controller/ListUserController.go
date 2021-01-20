package controller

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	OwnSQL "github.com/MathisBurger/OpenInventory-Backend/mysql"
	"github.com/gofiber/fiber/v2"
)

func ListUserController(c *fiber.Ctx) error {
	raw := string(c.Body())
	obj := models.LoginWithTokenRequest{}
	err := json.Unmarshal([]byte(raw), &obj)
	if err != nil {
		resp, _ := models.GetJsonResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(resp)
	}
	if OwnSQL.MySQL_loginWithToken(obj.Username, obj.Password, obj.Token) {
		conn := OwnSQL.GetConn()
		stmt, _ := conn.Prepare("SELECT `username`, `root`, `mail`, `register_date`, `status` FROM `inv_users`;")
		res, _ := stmt.Query()
		var answers []models.OutputUserStruct
		for res.Next() {
			var cache models.OutputUserStruct
			err = res.Scan(&cache.Username, &cache.Root, &cache.Mail, &cache.RegisterDate, &cache.Status)
			if err != nil {
				panic(err)
			}
			answers = append(answers, cache)
		}
		res.Close()
		stmt.Close()
		conn.Close()
		return c.JSON(models.ListUserResponseModel{
			Message: "successfully fetched user",
			Alert:   "alert alert-success",
			User:    answers,
		})
	} else {
		resp, _ := models.GetJsonResponse("You do not have the permission to perform this", "alert alert-danger", "ok", "None", 200)
		return c.Send(resp)
	}
}
