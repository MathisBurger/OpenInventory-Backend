package utils

import (
	"crypto/sha512"
	"encoding/base64"
)

// hashes with SHA512
func HashWithSalt(pwd string) string {
	hasher := sha512.New()
	hasher.Write([]byte(pwd))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha
}
