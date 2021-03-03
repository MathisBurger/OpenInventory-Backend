package models

// table model
type TableModel struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Entries    int    `json:"entries"`
	MinPermLvl int    `json:"min_perm_lvl"`
	CreatedAt  string `json:"created_at"`
}
