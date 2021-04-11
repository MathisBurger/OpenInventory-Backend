package models

import "database/sql"

type TwoFactorModel struct {
	ID     int    `json:"id"`
	Secret string `json:"secret"`
	Owner  string `json:"owner"`
}

func (c TwoFactorModel) Parse(resp *sql.Rows) TwoFactorModel {
	var mdl TwoFactorModel
	_ = resp.Scan(&mdl.ID, &mdl.Secret, &mdl.Owner)
	return mdl
}

func (c TwoFactorModel) ParseAll(resp *sql.Rows) []TwoFactorModel {
	var answers []TwoFactorModel

	for resp.Next() {
		var mdl TwoFactorModel
		_ = resp.Scan(&mdl.ID, &mdl.Secret, &mdl.Owner)
		answers = append(answers, mdl)
	}
	return answers
}
