package twoFactorAuth

import (
	"io/ioutil"
	"net/url"
	qr "rsc.io/qr"
)

// This function generates a QR code for the
// google authenticator, based on the account
// name and a random secret
func GenerateQR(acc string, secret string) {

	issuer := "OpenInventory"

	URL, err := url.Parse("otpauth://totp")
	if err != nil {
		panic(err)
	}

	URL.Path += "/" + url.PathEscape(issuer) + ":" + url.PathEscape(acc)

	params := url.Values{}
	params.Add("secret", secret)
	params.Add("issuer", issuer)

	URL.RawQuery = params.Encode()

	code, err := qr.Encode(URL.String(), qr.Q)
	if err != nil {
		panic(err.Error())
	}
	b := code.PNG()
	err = ioutil.WriteFile("./temp-qr/qr.png", b, 0600)
	if err != nil {
		panic(err.Error())
	}
}
