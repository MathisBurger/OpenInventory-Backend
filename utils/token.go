package utils

import (
	"math/rand"
	"strings"
	"time"
)

func GenerateToken() string {
	rand.Seed(time.Now().Unix())
	charSet := "abcdedfghijklmnopqrstuvwxyzABCDEFHGIJKLMNOPQRSTUVWXYZ0123456789#&%$§"
	var output strings.Builder
	length := 10
	for i := 0; i < length; i++ {
		random := rand.Intn(len(charSet))
		randomChar := charSet[random]
		output.WriteString(string(randomChar))
	}
	return output.String()
}
