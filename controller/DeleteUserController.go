package controller

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions/utils"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/gofiber/fiber/v2"
)

func DeleteUserController(c *fiber.Ctx) error {
	raw := string(c.Body())
	obj := models.DeleteUserRequestModel{}
	err := json.Unmarshal([]byte(raw), &obj)
	if err != nil {
		res, err := models.GetJSONResponse("Invaild JSON body", "alert alert-danger", "error", "None", 200)
		if err != nil {
			utils.LogError("[DeleteUserController.go, 18, InputError] " + err.Error())
		}
		return c.Send(res)
	}
	if !checkDeleteUserRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	status := actions.MysqlLoginWithTokenRoot(obj.Username, obj.Password, obj.Token)
	if status {
		conn := actions.GetConn()
		stmt, _ := conn.Prepare("DELETE FROM `inv_users` WHERE `username`=?;")
		res, _ := stmt.Exec(obj.User)
		aff, _ := res.RowsAffected()
		defer stmt.Close()
		defer conn.Close()
		if aff == 0 {
			resp, _ := models.GetJSONResponse("This user does not exist", "alert alert-warning", "ok", "None", 200)
			return c.Send(resp)
		}
		resp, _ := models.GetJSONResponse("Successfully deleted user", "alert alert-success", "ok", "None", 200)
		return c.Send(resp)
	}
	res, _ := models.GetJSONResponse("You do not have the permission to execute this", "alert alert-danger", "ok", "None", 200)
	return c.Send(res)

}

func checkDeleteUserRequest(obj models.DeleteUserRequestModel) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.User != ""
}
