package controller

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	OwnSQL "github.com/MathisBurger/OpenInventory-Backend/mysql"
	"github.com/gofiber/fiber/v2"
)

func DeleteUserController(c *fiber.Ctx) error {
	raw := string(c.Body())
	obj := models.DeleteUserRequestModel{}
	err := json.Unmarshal([]byte(raw), &obj)
	if err != nil {
		response, err := models.GetJsonResponse("Invaild JSON body", "alert alert-danger", "error", "None", 200)
		if err != nil {
			panic(err)
		}
		return c.Send(response)
	}
	status := OwnSQL.MySQL_loginWithToken_ROOT(obj.Username, obj.Password, obj.Token)
	if status {
		conn := OwnSQL.GetConn()
		stmt, _ := conn.Prepare("DELETE FROM `inv_users` WHERE `username`=?;")
		res, _ := stmt.Exec(obj.User)
		aff, _ := res.RowsAffected()
		stmt.Close()
		conn.Close()
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
