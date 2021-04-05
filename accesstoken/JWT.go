package accesstoken

import (
	"crypto/rsa"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"time"
)

// this defines the method for signing JWT token
var signMethod = jwt.SigningMethodRS256

// This implements the RSA keypair
// needed to sign JWT
type JWTKeyPair struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

// this functions generates a new Keypair for JWT signing
func NewJWTManager(privateKeyFile, publicKeyFile string) (m *JWTKeyPair, err error) {
	m = new(JWTKeyPair)

	var kd []byte
	if privateKeyFile != "" {
		kd, err = ioutil.ReadFile(privateKeyFile)
		if err != nil {
			return
		}
		m.privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(kd)
		if err != nil {
			return
		}
	}

	if publicKeyFile != "" {
		kd, err = ioutil.ReadFile(publicKeyFile)
		if err != nil {
			return
		}
		m.publicKey, err = jwt.ParseRSAPublicKeyFromPEM(kd)
		if err != nil {
			return
		}
	}

	return
}

func (m *JWTKeyPair) Generate(ident string, expire time.Duration) (token string, err error) {
	if m.privateKey == nil {
		err = errors.New("not supported with this instance")
		return
	}

	now := time.Now()
	token, err = jwt.NewWithClaims(signMethod, jwt.StandardClaims{
		Subject:   ident,
		ExpiresAt: now.Add(expire).Unix(),
		IssuedAt:  now.Unix(),
	}).SignedString(m.privateKey)

	return
}

func (m *JWTKeyPair) Validate(token string) (ident string, err error) {
	if m.publicKey == nil {
		err = errors.New("not supported with this instance")
		return
	}

	claims := new(jwt.StandardClaims)
	jwtToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return m.publicKey, nil
	})
	if err != nil {
		return
	}

	if !jwtToken.Valid {
		err = errors.New("invalid claims")
	}

	ident = claims.Subject

	return
}