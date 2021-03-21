package main

import (
	"fmt"
	config2 "github.com/MathisBurger/OpenInventory-Backend/config"
	"github.com/MathisBurger/OpenInventory-Backend/controller/general"
	"github.com/MathisBurger/OpenInventory-Backend/controller/permission-management"
	"github.com/MathisBurger/OpenInventory-Backend/controller/table-management"
	"github.com/MathisBurger/OpenInventory-Backend/controller/user-management"
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
		app.Get("/api", general.DefaultController)
		app.Get("/api/info", general.InformationController)
		app.Post("/api/login", general.LoginController)
		app.Post("/api/check-creds", general.CheckCredsController)

		// user management
		app.Get("/api/user-management/ListUser", user_management.ListUserController)
		app.Post("/api/user-management/AddUser", user_management.AddUserController)
		app.Delete("/api/user-management/DeleteUser", user_management.DeleteUserController)

		// table management
		app.Get("/api/table-management/getAllTables", table_management.GetAllTablesController)
		app.Post("/api/table-management/createTable", table_management.CreateTableController)
		app.Get("/api/table-management/getTableContent", table_management.GetTableContentController)
		app.Post("/api/table-management/AddTableEntry", table_management.AddTableEntryController)
		app.Post("/api/table-management/getTableColumns", table_management.GetTableColumnsController)
		app.Post("/api/table-management/RemoveTableEntry", table_management.RemoveTableEntryController)
		app.Post("/api/table-management/DeleteTable", table_management.DeleteTableController)
		app.Post("/api/table-management/editTableEntry", table_management.EditTableEntryController)
		app.Post("/api/table-management/renameTableColumn", table_management.RenameTableColumnController)
		app.Post("/api/table-management/renameTable", table_management.RenameTableController)

		// permission management
		app.Post("/api/permission-management/createPermissionGroup", permission_management.CreatePermissionGroupController)
		app.Post("api/permission-management/addUserToPermissionGroup", permission_management.AddUserToPermissionGroupController)
		app.Post("/api/permission-management/deletePermissionGroup", permission_management.DeletePermissionGroupController)
		app.Post("/api/permission-management/removeUserFromPermissionGroup", permission_management.RemoveUserFromPermissionGroupController)
		app.Post("/api/permission-management/editTableMinPermLvl", permission_management.EditTableMinPermLvlController)
		app.Post("/api/permission-management/listAllPermsOfUser", permission_management.ListAllPermOfUserController)
		app.Post("/api/permission-management/listAllPermGroupsOfTable", permission_management.ListAllPermGroupsOfTableController)
		app.Post("/api/permission-management/listAllPermissionGroups", permission_management.ListAllPermissionGroupsController)

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
