package entities

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Username string `gorm:"unique;not null;"`
	Email    string `gorm:"unique;not null;"`
	Password string `gorm:"not null" json:"-"`
	Bio      string

	LinkedInURL string
	GithubURL   string
	WebsiteURL  string

	IsAdmin    bool `gorm:"default:false;"`
	IsProposer bool `gorm:"default:false"`
	IsBanned   bool `gorm:"default:false"`

	UserIconURL string

	Problems    []Problem    `gorm:"foreignKey:AuthorId;references:ID;constraint:OnDelete:CASCADE;";json:"-"`
	Submissions []Submission `gorm:"foreignKey:UserId;references:ID;constraint:OnDelete:CASCADE;";json:"-"`
}
