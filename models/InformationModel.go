package models

import (
	"encoding/json"
	"runtime"
	"strconv"
)

type InformationModel map[string]string

// returns basic information
func GetInformationResponse() ([]byte, error) {
	msg := make(map[string]string)
	msg["api-version"] = "v0.0.3-dev"
	msg["api-language"] = "golang"
	msg["operating-system"] = runtime.GOOS
	msg["cpu-cores"] = strconv.Itoa(runtime.NumCPU())
	return json.MarshalIndent(msg, "", "  ")
}
