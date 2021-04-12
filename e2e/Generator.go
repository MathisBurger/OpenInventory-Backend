package e2e

import (
	"crypto/rand"
	"crypto/rsa"
	"log"
)

// This function generates RSA public and
// private keys for end to end encryption
func GenerateRSA_Keys(bits int) (*rsa.PrivateKey, *rsa.PublicKey) {
	privkey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		log.Println(err.Error())
	}
	return privkey, &privkey.PublicKey
}
