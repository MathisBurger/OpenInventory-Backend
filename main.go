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
	// check installation status
	if installation.Install() {
		config, _ := config2.ParseConfig()
		app := fiber.New(fiber.Config{
			Prefork: config.ServerCFG.Prefork,
		})

		// Logger configuration
		app.Use(logger.New())
		app.Use(cors.New())

		initWebpaths(app, config)

		// Basic GET Requests
		app.Get("/api", controller.DefaultController)
		app.Get("/api/info", controller.InformationController)

		// user management
		app.Post("/api/user-management/ListUser", controller.ListUserController)
		app.Post("/api/user-management/AddUser", controller.AddUserController)
		app.Post("/api/user-management/DeleteUser", controller.DeleteUserController)

		// table management
		app.Post("/api/login", controller.LoginController)
		app.Post("/api/check-creds", controller.CheckCredsController)
		app.Post("/api/table-management/getAllTables", controller.GetAllTablesController)
		app.Post("/api/table-management/createTable", controller.CreateTableController)
		app.Post("/api/table-management/getTableContent", controller.GetTableContentController)
		app.Post("/api/table-management/AddTableEntry", controller.AddTableEntryController)
		app.Post("/api/table-management/getTableColumns", controller.GetTableColumnsController)
		app.Post("/api/table-management/RemoveTableEntry", controller.RemoveTableEntryController)
		app.Post("/api/table-management/DeleteTable", controller.DeleteTableController)
		app.Post("/api/table-management/editTableEntry", controller.EditTableEntryController)
		app.Post("/api/table-management/renameTableColumn", controller.RenameTableColumnController)
		app.Post("/api/table-management/renameTable", controller.RenameTableController)

		// permission management
		app.Post("/api/permission-management/createPermissionGroup", controller.CreatePermissionGroupController)
		app.Post("api/permission-management/addUserToPermissionGroup", controller.AddUserToPermissionGroupController)
		app.Post("/api/permission-management/deletePermissionGroup", controller.DeletePermissionGroupController)
		app.Post("/api/permission-management/removeUserFromPermissionGroup", controller.RemoveUserFromPermissionGroupController)
		app.Post("/api/permission-management/editTableMinPermLvl", controller.EditTableMinPermLvlController)
		app.Post("/api/permission-management/listAllPermsOfUser", controller.ListAllPermOfUserController)
		app.Post("/api/permission-management/listAllPermGroupsOfTable", controller.ListAllPermGroupsOfTableController)
		app.Post("/api/permission-management/listAllPermissionGroups", controller.ListAllPermissionGroupsController)

		// App Configuration
		app.Listen(":" + config.ServerCFG.Port)
	} else {
		fmt.Println("Please fix errors first to launch webserver")
	}
}

func initWebpaths(app *fiber.App, cfg *config2.Config) {
	for k, v := range cfg.ServerCFG.WebEndpoints {
		app.Static(k, v)
		fmt.Println("initialized '" + k + "' as web-endpoint at '" + v + "'")
	}
}
