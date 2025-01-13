package models

import (
	"github.com/google/uuid"
)

type User struct {
	Id          uuid.UUID `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
	Name        string    `gorm:"not null" json:"name"`
	UrlTag      string    `gorm:"unique;not null;check:url_tag <> ''" json:"url_tag"`
	Description string    `json:"description"`
	ProfilePic  string    `json:"profile_pic"`
}

func New(id uuid.UUID, name string) *User {
	return &User{Id: id, Name: name, UrlTag: id.String(), Description: "", ProfilePic: ""}
}
