package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"
	"fmt"
)

// This function generates an RS256 key pair
// Only if files are not existing
func GenerateKeys() {
	if _, err := os.Stat("./certs/private.pem"); os.IsNotExist(err) {

		fmt.Println("JWT keys are missing")

		privKey, _ := rsa.GenerateKey(rand.Reader, 2048)
		privBytes := pem.EncodeToMemory(
			&pem.Block{
				Type:  "RSA PRIVATE KEY",
				Bytes: x509.MarshalPKCS1PrivateKey(privKey),
			},
		)
		ioutil.WriteFile("./certs/private.pem", privBytes, 0644)

		fmt.Println("Generated private key")

		public := &privKey.PublicKey

		pubASN1, _ := x509.MarshalPKIXPublicKey(public)
		pubBytes := pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubASN1,
		})
		ioutil.WriteFile("./certs/public.pem", pubBytes, 0644)

		fmt.Println("Generated public key from private key")
	} else {
		fmt.Println("JWT keys are existing")
	}
}