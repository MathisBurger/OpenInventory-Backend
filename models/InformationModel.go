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
	msg["api_version"] = "v0.0.4-production"
	msg["api_language"] = "golang"
	msg["operating_system"] = runtime.GOOS
	msg["cpu_cores"] = strconv.Itoa(runtime.NumCPU())
	return json.MarshalIndent(msg, "", "  ")
}
