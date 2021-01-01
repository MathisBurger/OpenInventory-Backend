package main

import (
	"fmt"
	config2 "github.com/MathisBurger/OpenInventory-Backend/config"
	"github.com/MathisBurger/OpenInventory-Backend/controller"
	"github.com/MathisBurger/OpenInventory-Backend/installation"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	if installation.Install() {
		config, _ := config2.ParseConfig()
		app := fiber.New(fiber.Config{
			Prefork: config.ServerCFG.Prefork,
		})

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
		app.Post("/table-management/createTable", controller.CreateTableController)
		app.Post("/table-management/getTableContent", controller.GetTableContentController)
		app.Post("/table-management/AddTableEntry", controller.AddTableEntryController)
		app.Post("/table-management/getTableColumns", controller.GetTableColumnsController)
		app.Post("/table-management/RemoveTableEntry", controller.RemoveTableEntryController)
		app.Post("/table-management/DeleteTable", controller.DeleteTableController)
		app.Post("/table-management/ListUser", controller.ListUserController)

		// App Configuration
		app.Listen(":" + config.ServerCFG.Port)
	} else {
		fmt.Println("Please fix errors first to launch webserver")
	}
}
