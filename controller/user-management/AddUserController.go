package user_management

import (
	"encoding/json"
	"github.com/MathisBurger/OpenInventory-Backend/e2e"
	"strings"

	"github.com/MathisBurger/OpenInventory-Backend/config"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/middleware"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

type addUserRequest struct {
	User struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Root     bool   `json:"root"`
		Mail     string `json:"mail"`
		Status   string `json:"status"`
	} `json:"user"`
}

/////////////////////////////////////////////////////////////
//                                                         //
//                    AddUserController                    //
//      This controller adds user to inventory system      //
//    It requires login credentials and information about  //
//                   the new user.                         //
//                                                         //
/////////////////////////////////////////////////////////////
func AddUserController(c *fiber.Ctx) error {

	// init and parse the request object
	obj := addUserRequest{}
	decrypted, err := e2e.DecryptBytes(c.Body())
	if err != nil {
		return c.SendStatus(400)
	}
	err = json.Unmarshal(decrypted, &obj)

	// check request
	if err != nil {
		if cfg, _ := config.ParseConfig(); cfg.ServerCFG.LogRequestErrors {
			utils.LogError(err.Error(), "AddUserController.go", 19)
		}

		res, _ := models.GetJSONResponse("Invaild JSON body", "#d41717", "error", "None", 200)
		return c.Send(res)
	}
	if !checkAddUserRequest(obj) {
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "#d41717", "ok", "None", 200)
		return c.Send(res)
	}

	// check if the username is too long
	if !checkUsernameLength(obj.User.Username) {
		res, _ := models.GetJSONResponse("This username is too long", "#d41717", "ok", "None", 200)
		return c.Send(res)
	}

	// check if the mail is valid
	if !checkEmail(obj.User.Mail) {
		res, _ := models.GetJSONResponse("Your mail is not a email-address", "#d41717", "ok", "None", 200)
		return c.Send(res)
	}

	conn := actions.GetConn()
	defer conn.Close()

	// check login status
	if ok, ident := middleware.ValidateAccessToken(c); ok && actions.CheckUserHasHigherPermission(conn, ident, 0, "default.root") {

		hash := utils.HashPassword(obj.User.Password)

		actions.AddUser(obj.User.Username, hash, obj.User.Root, obj.User.Mail, obj.User.Status)

		res, _ := models.GetJSONResponse("Successfully added user", "#1db004", "ok", "None", 200)
		return c.Send(res)
	}

	// login failed
	res, _ := models.GetJSONResponse("You do not have the permission to perform this", "#d41717", "ok", "None", 200)
	return c.Send(res)
}

// checks the request
// struct fields should not be default
func checkAddUserRequest(obj addUserRequest) bool {
	return obj.User != addUserRequest{}.User
}

// check if username is too long
func checkUsernameLength(username string) bool {
	split := strings.Split(username, "")
	return len(split) < 32
}

// check if email is valid
func checkEmail(mail string) bool {
	return strings.Contains(mail, "@") && len(strings.Split(mail, ".")) > 0
}
