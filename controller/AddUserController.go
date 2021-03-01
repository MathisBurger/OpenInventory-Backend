package controller

import (
	"github.com/MathisBurger/OpenInventory-Backend/config"
	"github.com/MathisBurger/OpenInventory-Backend/database/actions"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	"github.com/MathisBurger/OpenInventory-Backend/utils"
	"github.com/gofiber/fiber/v2"
	"strings"
)

// ---------------------------------------------
//             addUserRequest
//    This struct contains login credentials
//      and information about the new user
// ---------------------------------------------
type addUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
	User     struct {
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

	// initializing the request object
	obj := new(addUserRequest)

	// parsing the body into the request object
	err := c.BodyParser(obj)

	// returns "Wrong JSON syntax" response if error is unequal nil
	if err != nil {

		// checks if request errors should be logged
		if cfg, _ := config.ParseConfig(); cfg.ServerCFG.LogRequestErrors {

			// log error
			utils.LogError(err.Error(), "AddUserController.go", 19)
		}

		// returns response
		res, _ := models.GetJSONResponse("Invaild JSON body", "alert alert-danger", "error", "None", 200)
		return c.Send(res)
	}

	// check if request has been parsed correctly
	if !checkAddUserRequest(obj) {

		// returns "Wrong JSON syntax" response
		res, _ := models.GetJSONResponse("Wrong JSON syntax", "alert alert-danger", "ok", "None", 200)
		return c.Send(res)
	}

	// check if the username is too long
	if !checkUsernameLength(obj.User.Username) {

		// returns "This username is too long" response if username is too long
		res, _ := models.GetJSONResponse("This username is too long", "alert alert-warning", "ok", "None", 200)
		return c.Send(res)
	}

	// check if the mail is valid
	if !checkEmail(obj.User.Mail) {

		// returns "Your mail is not a email-address" if mail is not valid
		res, _ := models.GetJSONResponse("Your mail is not a email-address", "alert alert-warning", "ok", "None", 200)
		return c.Send(res)
	}

	// get status of login
	status := actions.MysqlLoginWithTokenRoot(obj.Username, obj.Password, obj.Token)

	// check login status
	if status {

		// hash the password
		hash := utils.HashWithSalt(obj.User.Password)

		// add user to database
		actions.AddUser(obj.User.Username, hash, obj.User.Root, obj.User.Mail, obj.User.Status)

		// returns "Successfully added user" response
		res, _ := models.GetJSONResponse("Successfully added user", "alert alert-success", "ok", "None", 200)
		return c.Send(res)
	}

	// returns "You do not have the permission to perform this" response if login failed
	res, _ := models.GetJSONResponse("You do not have the permission to perform this", "alert alert-danger", "ok", "None", 200)
	return c.Send(res)
}

//////////////////////////////////////////////////////////
//                                                      //
//                  checkAddUserRequest                 //
//              consumes the request object             //
//   checks if struct fields are not the default value  //
//                                                      //
//////////////////////////////////////////////////////////
func checkAddUserRequest(obj *addUserRequest) bool {
	return obj.Username != "" && obj.Password != "" && obj.Token != "" && obj.User != addUserRequest{}.User
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
