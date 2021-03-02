package models

// table model
type TableModel struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Entries    int    `json:"entries"`
	MinPermLvl int    `json:"min-perm-lvl"`
	CreatedAt  string `json:"created_at"`
}
