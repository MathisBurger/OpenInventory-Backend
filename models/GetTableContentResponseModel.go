package models

type GetTableContentResponseModel struct {
	Message    string `json:"message"`
	Alert      string `json:"alert"`
	Status     string `json:"status"`
	HttpStatus int    `json:"httpStatus"`
	Elements   string `json:"elements"`
}
