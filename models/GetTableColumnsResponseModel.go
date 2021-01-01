package models

type GetTableColumnsResponseModel struct {
	Message string      `json:"message"`
	Alert   string      `json:"alert"`
	Columns interface{} `json:"columns"`
}
