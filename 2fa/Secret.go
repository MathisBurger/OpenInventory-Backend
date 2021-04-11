package twoFactorAuth

import (
	"crypto/rand"
	"encoding/base32"
)

// This function generates a random
// secret for the 2fa authentication
func GenerateSecret() string {

	secret := make([]byte, 10)
	_, err := rand.Read(secret)
	if err != nil {
		panic(err)
	}

	return base32.StdEncoding.EncodeToString(secret)
}
