package controller

import (
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

type deleteUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
	User     string `json:"user"`
}

func DeleteUserController(c *fiber.Ctx) error {
	obj := new(deleteUserRequest)
	err := c.BodyParser(obj)
	if err != nil {
		res, err := models.GetJSONResponse("Invaild JSON body", "alert alert-danger", "error", "None", 200)
		if err != nil {
			utils.LogError(err.Error(), "DeleteUserController.go", 18)
		}
		return c.Send(res)
	}
	if !checkDeleteUserRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	status := actions.MysqlLoginWithTokenRoot(obj.Username, obj.Password, obj.Token)
	if status {
		actions.DeleteUser(obj.User)
		resp, _ := models.GetJSONResponse("Successfully deleted user", "alert alert-success", "ok", "None", 200)
		return c.Send(resp)
	}
	res, _ := models.GetJSONResponse("You do not have the permission to execute this", "alert alert-danger", "ok", "None", 200)
	return c.Send(res)

}

func checkDeleteUserRequest(obj *deleteUserRequest) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.User != ""
}
