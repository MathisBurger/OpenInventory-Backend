package models

import (
	"database/sql"
	"time"
)

// global model
type UserModel struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Password     string    `json:"password"`
	Token        string    `json:"token"`
	Permissions  string    `json:"permissions"`
	Root         bool      `json:"root"`
	Mail         string    `json:"mail"`
	Displayname  string    `json:"displayname"`
	RegisterDate time.Time `json:"register_date"`
	Status       string    `json:"status"`
}

// fetch only one [resp.Next() required]
func (c UserModel) Parse(resp *sql.Rows) UserModel {
	var mdl UserModel
	_ = resp.Scan(&mdl.ID, &mdl.Username, &mdl.Password, &mdl.Token, &mdl.Permissions, &mdl.Root, &mdl.Mail, &mdl.Displayname, &mdl.RegisterDate, &mdl.Status)
	return mdl
}

// fetch all
func (c UserModel) ParseAll(resp *sql.Rows) []UserModel {
	var answers []UserModel
	for resp.Next() {
		var mdl UserModel
		_ = resp.Scan(&mdl.ID, &mdl.Username, &mdl.Password, &mdl.Token, &mdl.Permissions, &mdl.Root, &mdl.Mail, &mdl.Displayname, &mdl.RegisterDate, &mdl.Status)
		answers = append(answers, mdl)
	}
	return answers
}
