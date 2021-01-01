package models

type RemoveTableEntryRequestModel struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Token     string `json:"token"`
	TableName string `json:"table_name"`
	RowID     int    `json:"row_id"`
}
