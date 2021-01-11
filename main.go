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

		// Static Web Files
		app.Static("/", "./web")
		app.Static("/login", "./web/index.html")
		app.Static("/dashboard", "./web/index.html")

		// Basic GET Requests
		app.Get("/api", controller.DefaultController)
		app.Get("/api/info", controller.InformationController)

		// POST Requests
		app.Post("/api/login", controller.LoginController)
		app.Post("/api/check-creds", controller.CheckCredsController)
		app.Post("/api/table-management/getAllTables", controller.GetAllTablesController)
		app.Post("/api/table-management/createTable", controller.CreateTableController)
		app.Post("/api/table-management/getTableContent", controller.GetTableContentController)
		app.Post("/api/table-management/AddTableEntry", controller.AddTableEntryController)
		app.Post("/api/table-management/getTableColumns", controller.GetTableColumnsController)
		app.Post("/api/table-management/RemoveTableEntry", controller.RemoveTableEntryController)
		app.Post("/api/table-management/DeleteTable", controller.DeleteTableController)
		app.Post("/api/table-management/ListUser", controller.ListUserController)
		app.Post("/api/table-management/AddUser", controller.AddUserController)
		app.Post("/api/table-management/DeleteUser", controller.DeleteUserController)

		// App Configuration
		app.Listen(":" + config.ServerCFG.Port)
	} else {
		fmt.Println("Please fix errors first to launch webserver")
	}
}
