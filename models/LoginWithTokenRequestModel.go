package models

// login with token request
type LoginWithTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
}
