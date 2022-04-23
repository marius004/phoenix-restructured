package entities

import "gorm.io/gorm"

type Problem struct {
	gorm.Model

	Name        string `gorm:"unique"`
	Description string
	Difficulty  string

	Visible  bool
	AuthorId uint

	TimeLimit   float32
	MemoryLimit int
	StackLimit  int

	ProblemTests []ProblemTest `gorm:"foreignKey:ProblemId;references:ID;constraint:OnDelete:CASCADE;";json:"-"`
	Submissions  []Submission  `gorm:"foreignKey:ProblemId;references:ID;constraint:OnDelete:CASCADE;";json:"-"`
}
