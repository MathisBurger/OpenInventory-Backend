package models

type TableModel struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Entries   int    `json:"entries"`
	CreatedAt string `json:"created_at"`
}
