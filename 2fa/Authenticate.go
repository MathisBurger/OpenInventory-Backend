package twoFactorAuth

import (
	dgoogauth "github.com/dgryski/dgoogauth"
)

// This function needs the secret and the 2fa code
// It authorizes the code
func Authenticate(secret string, token string) bool {

	otpc := &dgoogauth.OTPConfig{
		Secret:      secret,
		WindowSize:  3,
		HotpCounter: 0,
	}

	val, err := otpc.Authenticate(token)
	if err != nil {
		return false
	}

	return val
}
