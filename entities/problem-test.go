package entities

import "gorm.io/gorm"

type ProblemTest struct {
	gorm.Model

	ProblemId uint
	Score     int `gorm:"default:10"`

	Input  []byte `gorm:"default:null"`
	Output []byte `gorm:"default:null"`

	SubmissionTests []SubmissionTest `gorm:"foreignKey:ProblemTestId;constraint:OnDelete:CASCADE;references:ID";json:"-"`
}

// TODO check for batch delete https://gorm.io/docs/delete.html
