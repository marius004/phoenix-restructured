package entities

import "gorm.io/gorm"

type SubmissionTest struct {
	gorm.Model

	ExecutionMessage string
	Score            int
	Time             float64
	Memory           int
	ExitCode         int `gorm:"default:0"`

	SubmissionId  uint
	ProblemTestId uint
}
