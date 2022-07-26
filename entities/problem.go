package entities

import "gorm.io/gorm"

const (
	Easy    = "easy"
	Medium  = "medium"
	Hard    = "hard"
	Contest = "contest"
)

var Difficulties = []string{Easy, Medium, Hard, Contest}

type ProblemStatus string

func (status ProblemStatus) IsValid() bool {
	return status == UnPublished ||
		status == WaitingForApproval ||
		status == Published
}

const (
	UnPublished        = ProblemStatus("unpublished")
	WaitingForApproval = ProblemStatus("waiting for approval")
	Published          = ProblemStatus("published")
)

type Problem struct {
	gorm.Model

	Name        string `gorm:"unique"`
	Description string
	Difficulty  string `gorm:"default:easy"`

	Status   ProblemStatus `gorm:"default:unpublished"`
	AuthorId uint

	TimeLimit   float32
	MemoryLimit int
	StackLimit  int

	ProblemTests []ProblemTest `gorm:"foreignKey:ProblemId;references:ID;constraint:OnDelete:CASCADE;";json:"-"`
	Submissions  []Submission  `gorm:"foreignKey:ProblemId;references:ID;constraint:OnDelete:CASCADE;";json:"-"`
}
