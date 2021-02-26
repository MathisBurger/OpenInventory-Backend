package models

import "database/sql"

type TableModel struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Entrys     int    `json:"entrys"`
	MinPermLvl int    `json:"min_perm_lvl"`
	CreatedAt  string `json:"created_at"`
}

func (c TableModel) Parse(resp *sql.Rows) TableModel {
	var mdl TableModel
	_ = resp.Scan(&mdl.ID, &mdl.Name, &mdl.Entrys, &mdl.MinPermLvl, &mdl.CreatedAt)
	return mdl
}

func (c TableModel) ParseAll(resp *sql.Rows) []TableModel {
	var answers []TableModel
	for resp.Next() {
		var mdl TableModel
		_ = resp.Scan(&mdl.ID, &mdl.Name, &mdl.Entrys, &mdl.MinPermLvl, &mdl.CreatedAt)
		answers = append(answers, mdl)
	}
	return answers
}
