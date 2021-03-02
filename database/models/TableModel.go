package models

import "database/sql"

// global model
type TableModel struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Entrys     int    `json:"entrys"`
	MinPermLvl int    `json:"min_perm_lvl"`
	CreatedAt  string `json:"created_at"`
}

// fetch only one [resp.Next() required]
func (c TableModel) Parse(resp *sql.Rows) TableModel {
	var mdl TableModel
	_ = resp.Scan(&mdl.ID, &mdl.Name, &mdl.Entrys, &mdl.MinPermLvl, &mdl.CreatedAt)
	return mdl
}

// fetch all
func (c TableModel) ParseAll(resp *sql.Rows) []TableModel {
	var answers []TableModel
	for resp.Next() {
		var mdl TableModel
		_ = resp.Scan(&mdl.ID, &mdl.Name, &mdl.Entrys, &mdl.MinPermLvl, &mdl.CreatedAt)
		answers = append(answers, mdl)
	}
	return answers
}
