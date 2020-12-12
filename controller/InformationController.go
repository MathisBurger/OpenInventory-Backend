package controller

import (
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/gofiber/fiber/v2"
)

func InformationController(c *fiber.Ctx) error {
	response, err := models.GetInformationResponse()
	if err != nil {
		panic(err)
	}
	return c.Send(response)
}
