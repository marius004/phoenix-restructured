package entities

import "gorm.io/gorm"

// Submission Status
type SubmissionStatus string

const (
	Waiting    = SubmissionStatus("waiting")
	Evaluating = SubmissionStatus("evaluating")
	Evaluated  = SubmissionStatus("evaluated")
)

func (status SubmissionStatus) IsWaiting() bool    { return status == Waiting }
func (status SubmissionStatus) IsEvaluating() bool { return status == Evaluating }
func (status SubmissionStatus) IsEvaluated() bool  { return status == Evaluated }

// Programming Language
type ProgrammingLanguage string

const CPP = ProgrammingLanguage("c++")

func (lang ProgrammingLanguage) IsCPP() bool { return lang == CPP }

type Submission struct {
	gorm.Model

	Score      int                 `gorm:"default:0"`
	Language   ProgrammingLanguage `gorm:"default:"c++""`
	SourceCode []byte

	Status              SubmissionStatus `gorm:"default:"waiting""`
	CompilationMessage  string           `gorm:"default:"""`
	CompiledSuccesfully *bool            `gorm:"default:"null""`

	UserId    uint
	ProblemId uint
}
