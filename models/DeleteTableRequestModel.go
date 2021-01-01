package models

type DeleteTableRequestModel struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Token     string `json:"token"`
	TableName string `json:"table_name"`
}
