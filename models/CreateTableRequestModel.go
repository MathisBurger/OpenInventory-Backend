package models

type CreateTableRequestModel struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Token     string `json:"token"`
	TableName string `json:"table_name"`
	RowConfig string `json:"row_config"`
}
