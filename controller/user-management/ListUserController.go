package user_management

import (
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/middleware"
	"github.com/MathisBurger/OpenInventory-Backend/models"
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

	// check login
	if ok, _ := middleware.ValidateAccessToken(c); ok {

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
