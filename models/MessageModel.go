package models

import (
	"encoding/json"
	"strconv"
)

type MessageModel map[string]string

func GetJSONResponse(
	message string,
	alert string,
	status string,
	token string,
	httpStatus int,
) ([]byte, error) {
	msg := make(map[string]string)
	msg["message"] = message
	msg["alert"] = alert
	msg["status"] = status
	msg["token"] = token
	msg["httpStatus"] = strconv.Itoa(httpStatus)
	return json.MarshalIndent(msg, "", "  ")
}
