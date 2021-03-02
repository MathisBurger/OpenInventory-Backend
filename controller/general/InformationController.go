package general

import (
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/gofiber/fiber/v2"
)

////////////////////////////////////////////////////////////////////
//                                                                //
//                EditTableMinPermLvlController                   //
//           This controller returns basic information            //
//                                                                //
////////////////////////////////////////////////////////////////////
func InformationController(c *fiber.Ctx) error {
	res, err := models.GetInformationResponse()
	if err != nil {
		panic(err)
	}
	return c.Send(res)
}
