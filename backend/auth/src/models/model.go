package models

import (
	"github.com/google/uuid"
)

type User struct {
	Id    uuid.UUID `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
	Login string    `gorm:"unique;not null;check:login <> ''"`
	Name  string    `gorm:"not null;check:name <> ''"`
	Pass  string    `gorm:"not null;check:pass <> ''"`
}

func New(login string, name string, pass string) (*User, error) {
	user := User{
		Id:    uuid.New(),
		Login: login,
		Name:  name,
		Pass:  pass,
	}

	return &user, nil
}
