package e2e

import (
	"io/ioutil"
	"os"
)

// This function generates a new key pair if not exists
// and saves them into pem files
func SaveKeys() {

	if _, err := os.Stat("./certs/e2e-private.pem"); os.IsNotExist(err) {

		priv, pub := GenerateRSA_Keys(2048)

		priv_bytes := PrivateKeyToBytes(priv)
		ioutil.WriteFile("./certs/e2e-private.pem", priv_bytes, 0644)

		pub_bytes := PublicKeyToBytes(pub)
		ioutil.WriteFile("./certs/e2e-public.pem", pub_bytes, 0644)
	}
}
