package models

type ListUserResponseModel struct {
	Message string             `json:"message"`
	Alert   string             `json:"alert"`
	User    []OutputUserStruct `json:"user"`
}

type OutputUserStruct struct {
	Username     string  `json:"username"`
	Root         bool    `json:"root"`
	Mail         string  `json:"mail"`
	RegisterDate []uint8 `json:"register_date"`
	Status       string  `json:"status"`
}
