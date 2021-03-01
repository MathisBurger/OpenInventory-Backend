package controller

import (
	"github.com/MathisBurger/OpenInventory-Backend/config"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

type listUserResponse struct {
	Message string                    `json:"message"`
	Alert   string                    `json:"alert"`
	User    []models.OutputUserStruct `json:"user"`
}

func ListUserController(c *fiber.Ctx) error {
	obj := new(models.LoginWithTokenRequest)
	err := c.BodyParser(obj)
	if err != nil {
		if cfg, _ := config.ParseConfig(); cfg.ServerCFG.LogRequestErrors {
			utils.LogError(err.Error(), "ListUserController.go", 16)
		}
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if !checkCheckCredsRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}
	if actions.MysqlLoginWithToken(obj.Username, obj.Password, obj.Token) {
		user := actions.GetAllUser()
		var outputUser []models.OutputUserStruct
		for _, el := range user {
			outputUser = append(outputUser, models.OutputUserStruct{
				el.Username,
				el.Root,
				el.Mail,
				el.RegisterDate,
				el.Status,
			})
		}
		return c.JSON(listUserResponse{
			Message: "successfully fetched user",
			Alert:   "alert alert-success",
			User:    outputUser,
		})
	}
	res, _ := models.GetJSONResponse("You do not have the permission to perform this", "alert alert-danger", "ok", "None", 200)
	return c.Send(res)

}
