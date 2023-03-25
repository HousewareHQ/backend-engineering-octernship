package users

import (
	organisations "github.com/HousewareHQ/backend-engineering-octernship/test/_your_app_/models/organisation"
)

type User struct {
	Id           uint                       `json:"id" gorm:"primaryKey"`
	Username     string                     `json:"username" gorm:"unique"` // ! unique
	Password     []byte                     `json:"-"`                      // ! don't return password
	IsAdmin      bool                       `json:"is_admin" gorm:"default:false"`
	OrgId        uint                       `json:"org_id"`
	Organisation organisations.Organisation `json:"-" gorm:"foreignKey:OrgId"`
}
