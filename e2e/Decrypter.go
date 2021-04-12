package e2e

import (
	"crypto/rand"
	"crypto/rsa"
)

// DecryptWithPrivateKey decrypts data with private key
func DecryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) ([]byte, error) {

	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}
