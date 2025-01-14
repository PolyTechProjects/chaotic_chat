package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type Chat struct {
	Id          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primary_key"`
	Name        string    `gorm:"not null;check:name <> ''"`
	CreatorId   uuid.UUID `gorm:"type:uuid"`
	IsChannel   bool      `gorm:"not null;default:false"`
	JoinLink    string    `gorm:"unique;not null;check:join_link <> ''"`
	Description string
	ProfilePic  string
}

type UserChat struct {
	gorm.Model
	ChatId   uuid.UUID `gorm:"type:uuid"`
	UserId   uuid.UUID `gorm:"type:uuid"`
	ReadOnly bool
	IsAdmin  bool
}
