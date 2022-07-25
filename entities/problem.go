package entities

import "gorm.io/gorm"

// problem difficulty
const (
	Easy    = "easy"
	Medium  = "medium"
	Hard    = "hard"
	Contest = "contest"
)

type ProblemStatus string

const (
	UnPublished        = ProblemStatus("unpublished")
	WaitingForApproval = ProblemStatus("waiting for approval")
	Published          = ProblemStatus("published")
)

type Problem struct {
	gorm.Model

	Name        string `gorm:"unique"`
	Description string
	Difficulty  string

	Status   ProblemStatus `gorm:"default:unpublished"`
	AuthorId uint

	TimeLimit   float32
	MemoryLimit int
	StackLimit  int

	ProblemTests []ProblemTest `gorm:"foreignKey:ProblemId;references:ID;constraint:OnDelete:CASCADE;";json:"-"`
	Submissions  []Submission  `gorm:"foreignKey:ProblemId;references:ID;constraint:OnDelete:CASCADE;";json:"-"`
}
