package utils

import (
	"encoding/base64"
	"math/rand"
)


func ByteArray(lngth int) ([]byte, error) {
	arr := make([]byte, lngth)
	_, err := rand.Read(arr)
	return arr, err
}

func Base64(lngth int) string {
	str, _ := ByteArray(lngth)
	return base64.StdEncoding.EncodeToString(str)
}