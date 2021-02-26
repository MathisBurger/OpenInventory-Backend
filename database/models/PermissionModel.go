package models

import "database/sql"

type PermissionModel struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Color           string `json:"color"`
	PermissionLevel int    `json:"permission_level"`
}

func (c PermissionModel) Parse(resp *sql.Rows) PermissionModel {
	var mdl PermissionModel
	_ = resp.Scan(&mdl.ID, &mdl.Name, &mdl.Color, &mdl.PermissionLevel)
	return mdl
}