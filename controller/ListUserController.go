package controller

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/config"
	"github.com/MathisBurger/OpenInventory-Backend/controller/general"
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

////////////////////////////////////////////////////////////////////
//                                                                //
//                        ListUserController                      //
//                 This controller fetches all user               //
//         It requires models.LoginWithTokenRequest instance      //
//                                                                //
////////////////////////////////////////////////////////////////////
func ListUserController(c *fiber.Ctx) error {

	// init and parse the request object
	obj := models.LoginWithTokenRequest{}
	err := json.Unmarshal(c.Body(), &obj)

	// check request
	if err != nil {
		if cfg, _ := config.ParseConfig(); cfg.ServerCFG.LogRequestErrors {
			utils.LogError(err.Error(), "ListUserController.go", 16)
		}
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "#d41717", "ok", "None", 200)
		return c.Send(res)
	}
	if !general.CheckCheckCredsRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "#d41717", "ok", "None", 200)
		return c.Send(res)
	}

	// check login
	if actions.MysqlLoginWithToken(obj.Username, obj.Password, obj.Token) {

		user := actions.GetAllUser()
		var outputUser []models.OutputUserStruct

		// parse all user to output user
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
			Alert:   "#1db004",
			User:    outputUser,
		})
	}

	// no permission
	res, _ := models.GetJSONResponse("You do not have the permission to perform this", "#d41717", "ok", "None", 200)
	return c.Send(res)
}
