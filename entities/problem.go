package entities

import "gorm.io/gorm"

type Problem struct {
	gorm.Model

	Name        string `gorm:"unique"`
	Description string

	Visible  bool
	AuthorId uint

	TimeLimit   float32
	MemoryLimit int
	StackLimit  int

	ProblemTests []ProblemTest `gorm:"foreignKey:ProblemId;references:ID";json:"-"`
	Submissions  []Submission  `gorm:"foreignKey:ProblemId;references:ID";json:"-"`
}
