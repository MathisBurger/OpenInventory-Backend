package models

// login request struct
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
