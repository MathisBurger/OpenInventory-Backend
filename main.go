package main

import (
	"fmt"
	"github.com/MathisBurger/OpenInventory-Backend/controller"
	"github.com/MathisBurger/OpenInventory-Backend/installation"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	if installation.Install() {
		app := fiber.New()

		// Logger configuration
		app.Use(logger.New())

		// Web
		app.Static("/", "./web")

		// Basic GET Requests
		app.Get("/api", controller.DefaultController)
		app.Get("/api/info", controller.InformationController)

		// POST Requests
		app.Post("/api/login", controller.LoginController)
		app.Post("/api/check-creds", controller.CheckCredsController)

		// App Configuration
		app.Listen(":8080")
	} else {
		fmt.Println("Please fix errors first to launch webserver")
	}
}
