package controller

import (
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/gofiber/fiber/v2"
)

func DefaultController(c *fiber.Ctx) error {
	response, err := models.GetJsonResponse("API online", "alert alert-success", "ok", "None", 200)
	if err != nil {
		panic(err)
	}
	return c.Send(response)

}
