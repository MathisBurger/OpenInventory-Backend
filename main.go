package main

import (
	"fmt"
	"github.com/MathisBurger/OpenInventory-Backend/auth"
	config2 "github.com/MathisBurger/OpenInventory-Backend/config"
	"github.com/MathisBurger/OpenInventory-Backend/controller/general"
	permission_management "github.com/MathisBurger/OpenInventory-Backend/controller/permission-management"
	table_management "github.com/MathisBurger/OpenInventory-Backend/controller/table-management"
	user_management "github.com/MathisBurger/OpenInventory-Backend/controller/user-management"
	"github.com/MathisBurger/OpenInventory-Backend/e2e"
	"github.com/MathisBurger/OpenInventory-Backend/installation"
	"github.com/MathisBurger/OpenInventory-Backend/middleware"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"io/ioutil"
)

func main() {

	ioutil.WriteFile("./VERSION", []byte("v0.0.5-production"), 0644)

	// check installation status
	if installation.Install() {

		utils.GenerateKeys()

		e2e.SaveKeys()

		middleware.TwoFactorService()

		config, _ := config2.ParseConfig()

		app := fiber.New(fiber.Config{
			Prefork: config.ServerCFG.Prefork,
		})

		// Logger configuration
		app.Use(logger.New())
		app.Use(cors.New(cors.Config{
			AllowCredentials: true,
			ExposeHeaders:    "Authorization",
			AllowOrigins:     "http://127.0.0.1:4200",
		}))

		initWebpaths(app, config)

		// Basic GET Requests
		app.Get("/api", general.DefaultController)
		app.Get("/api/info", general.InformationController)
		app.Get("/api/publicKey", general.GetPublicKeyController)

		// Auth Endpoints
		app.Post("/api/auth/login", auth.LoginController)
		app.Get("/api/auth/accessToken", auth.AccessTokenController)
		app.Post("/api/auth/revokeSession", auth.RevokeSessionController)
		app.Get("/api/auth/me", auth.StatusController)
		app.Post("/api/auth/2fa", auth.TwoFactorAuthController)

		// user management
		app.Get("/api/user-management/ListUser", user_management.ListUserController)
		app.Post("/api/user-management/AddUser", user_management.AddUserController)
		app.Post("/api/user-management/DeleteUser", user_management.DeleteUserController)
		app.Patch("/api/user-management/Enable2FA", user_management.EnableTwoFactorController)

		// table management
		app.Get("/api/table-management/getAllTables", table_management.GetAllTablesController)
		app.Post("/api/table-management/createTable", table_management.CreateTableController)
		app.Get("/api/table-management/getTableContent", table_management.GetTableContentController)
		app.Post("/api/table-management/AddTableEntry", table_management.AddTableEntryController)
		app.Get("/api/table-management/getTableColumns", table_management.GetTableColumnsController)
		app.Post("/api/table-management/RemoveTableEntry", table_management.RemoveTableEntryController)
		app.Post("/api/table-management/DeleteTable", table_management.DeleteTableController)
		app.Patch("/api/table-management/editTableEntry", table_management.EditTableEntryController)
		app.Patch("/api/table-management/renameTableColumn", table_management.RenameTableColumnController)
		app.Patch("/api/table-management/renameTable", table_management.RenameTableController)

		// permission management
		app.Post("/api/permission-management/createPermissionGroup", permission_management.CreatePermissionGroupController)
		app.Post("api/permission-management/addUserToPermissionGroup", permission_management.AddUserToPermissionGroupController)
		app.Post("/api/permission-management/deletePermissionGroup", permission_management.DeletePermissionGroupController)
		app.Post("/api/permission-management/removeUserFromPermissionGroup", permission_management.RemoveUserFromPermissionGroupController)
		app.Patch("/api/permission-management/editTableMinPermLvl", permission_management.EditTableMinPermLvlController)
		app.Get("/api/permission-management/listAllPermsOfUser", permission_management.ListAllPermOfUserController)
		app.Get("/api/permission-management/listAllPermGroupsOfTable", permission_management.ListAllPermGroupsOfTableController)
		app.Get("/api/permission-management/listAllPermissionGroups", permission_management.ListAllPermissionGroupsController)

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
