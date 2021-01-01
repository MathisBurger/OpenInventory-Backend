package models

type AddUserRequestModel struct {
	Username string        `json:"username"`
	Password string        `json:"password"`
	Token    string        `json:"token"`
	User     AddUserStruct `json:"user"`
}

type AddUserStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Root     bool   `json:"root"`
	Mail     string `json:"mail"`
	Status   string `json:"status"`
}
