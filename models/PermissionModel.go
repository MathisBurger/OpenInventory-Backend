package models

// permission model
type PermissionModel struct {
	ID              int    `json:"ID"`
	Name            string `json:"name"`
	Color           string `json:"color"`
	PermissionLevel int    `json:"permission-level"`
}
