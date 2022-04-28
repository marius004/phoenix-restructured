package entities

import "gorm.io/gorm"

type Post struct {
	gorm.Model

	Title   string `gorm:"unique;not null;"`
	Content []byte `gorm:"not null;"`

	UserId uint
}
