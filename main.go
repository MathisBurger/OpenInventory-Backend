package main

import (
	"fmt"
	"github.com/MathisBurger/OpenInventory-Backend/controller"
	"github.com/MathisBurger/OpenInventory-Backend/installation"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	if installation.Install() {
		app := fiber.New()

		// Logger configuration
		app.Use(logger.New())
		app.Use(cors.New())

		// Basic GET Requests
		app.Get("/", controller.DefaultController)
		app.Get("/info", controller.InformationController)

		// POST Requests
		app.Post("/login", controller.LoginController)
		app.Post("/check-creds", controller.CheckCredsController)
		app.Post("/table-management/getAllTables", controller.GetAllTablesController)

		// App Configuration
		app.Listen(":8080")
	} else {
		fmt.Println("Please fix errors first to launch webserver")
	}
}
