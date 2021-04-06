package utils

import (
	"github.com/alexedwards/argon2id"
)

var params = &argon2id.Params{
	Memory:      128 * 1024,
	Iterations:  5,
	Parallelism: 5,
	SaltLength:  16,
	KeyLength:   32,
}

// hashes with SHA512
func HashPassword(pwd string) string {

	hash, err := argon2id.CreateHash(pwd, params)
	if err != nil {
		LogError(err.Error(), "hashing.go", 18)
	}
	return hash
}

func ValidateHash(pwd string, hash string) bool {

	match, err := argon2id.ComparePasswordAndHash(pwd, hash)
	if err != nil {
		LogError(err.Error(), "hashing.go", 24)
	}
	return match
}
