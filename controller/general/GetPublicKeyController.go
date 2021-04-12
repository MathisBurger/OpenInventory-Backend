package general

import (
	"github.com/gofiber/fiber/v2"
	"io/ioutil"
)

// This controller returns the RSA public key
func GetPublicKeyController(c *fiber.Ctx) error {

	bytes, _ := ioutil.ReadFile("./certs/e2e-public.pem")
	return c.Send(bytes)
}
