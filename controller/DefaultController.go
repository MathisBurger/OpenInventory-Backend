package controller

import (
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/gofiber/fiber/v2"
)

/////////////////////////////////////////////////////////////
//                                                         //
//                   DefaultController                     //
//          This controller returns a static response      //
//                                                         //
/////////////////////////////////////////////////////////////
func DefaultController(c *fiber.Ctx) error {

	res, _ := models.GetJSONResponse("API online", "alert alert-success", "ok", "None", 200)
	return c.Send(res)

}
