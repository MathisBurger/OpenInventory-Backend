package e2e

import (
	"crypto/rsa"
	"crypto/x509"
	b64 "encoding/base64"
	"encoding/pem"
	"io/ioutil"
	"log"
)

// This function parses the private
// key to bytes
func PrivateKeyToBytes(priv *rsa.PrivateKey) []byte {
	privBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		},
	)

	return privBytes
}

// PublicKeyToBytes public key to bytes
func PublicKeyToBytes(pub *rsa.PublicKey) []byte {
	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		log.Fatal(err)
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})

	return pubBytes
}

// DecryptBytes decrypts json body to bytes
func DecryptBytes(encrypted []byte) ([]byte, error) {
	rawKey, err := ioutil.ReadFile("./certs/e2e-private.pem")
	if err != nil {
		return nil, err
	}
	key := BytesToPrivateKey(rawKey)
	decoded, err := b64.StdEncoding.DecodeString(string(encrypted))
	if err != nil {
		return nil, err
	}
	decrypted, err := DecryptWithPrivateKey(decoded, key)
	if err != nil {
		return nil, err
	}
	return decrypted, nil
}

// BytesToPrivateKey bytes to private key
func BytesToPrivateKey(priv []byte) *rsa.PrivateKey {
	block, _ := pem.Decode(priv)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
	key, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		log.Fatal(err)
	}
	return key
}
