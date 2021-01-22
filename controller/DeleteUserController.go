package controller

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	OwnSQL "github.com/MathisBurger/OpenInventory-Backend/mysql"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

func DeleteUserController(c *fiber.Ctx) error {
	raw := string(c.Body())
	obj := models.DeleteUserRequestModel{}
	err := json.Unmarshal([]byte(raw), &obj)
	if err != nil {
		response, err := models.GetJsonResponse("Invaild JSON body", "alert alert-danger", "error", "None", 200)
		if err != nil {
			utils.LogError("[DeleteUserController.go, 18, InputError] " + err.Error())
		}
		return c.Send(response)
	}
	if !checkDeleteUserRequest(obj) {
		resp, _ := models.GetJsonResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(resp)
	}
	status := OwnSQL.MySQL_loginWithToken_ROOT(obj.Username, obj.Password, obj.Token)
	if status {
		conn := OwnSQL.GetConn()
		stmt, _ := conn.Prepare("DELETE FROM `inv_users` WHERE `username`=?;")
		res, _ := stmt.Exec(obj.User)
		aff, _ := res.RowsAffected()
		defer stmt.Close()
		defer conn.Close()
		if aff == 0 {
			resp, _ := models.GetJsonResponse("This user does not exist", "alert alert-warning", "ok", "None", 200)
			return c.Send(resp)
		} else {
			resp, _ := models.GetJsonResponse("Successfully deleted user", "alert alert-success", "ok", "None", 200)
			return c.Send(resp)
		}
	} else {
		resp, _ := models.GetJsonResponse("You do not have the permission to execute this", "alert alert-danger", "ok", "None", 200)
		return c.Send(resp)
	}
}

func checkDeleteUserRequest(obj models.DeleteUserRequestModel) bool {
	if obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.User != "" {
		return true
	} else {
		return false
	}
}
