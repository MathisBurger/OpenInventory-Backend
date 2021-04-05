
package models

import (
	"database/sql"
	"time"
)

// struct defines database model of refresh token
type RefreshTokenModel struct {
	ID       int
	Username string
	Token    string
	Deadline time.Time
}

func (c RefreshTokenModel) Parse(resp *sql.Rows) (bool, RefreshTokenModel) {
	if !resp.Next() {
		return false, RefreshTokenModel{}
	}
	var mdl RefreshTokenModel
	_ = resp.Scan(&mdl.ID, &mdl.Username, &mdl.Token, &mdl.Deadline)
	return true, mdl
}