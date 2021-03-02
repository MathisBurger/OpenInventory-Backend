package utils

import (
	"math/rand"
	"strings"
	"time"
)

// generates token
func GenerateToken() string {
	rand.Seed(time.Now().Unix())

	// defined charset
	charSet := "abcdedfghijklmnopqrstuvwxyzABCDEFHGIJKLMNOPQRSTUVWXYZ0123456789"

	var output strings.Builder
	length := 32

	// adding 32 random chars of charset
	for i := 0; i < length; i++ {
		random := rand.Intn(len(charSet))
		randomChar := charSet[random]
		output.WriteString(string(randomChar))
	}
	return output.String()
}
