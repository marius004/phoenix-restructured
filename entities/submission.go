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

const C = ProgrammingLanguage("c")

func (lang ProgrammingLanguage) IsC() bool { return lang == C }

type Submission struct {
	gorm.Model

	Score      int                 `gorm:"default:0"`
	Language   ProgrammingLanguage `gorm:"default:c"`
	SourceCode []byte

	Status              SubmissionStatus `gorm:"default:waiting"`
	CompiledSuccesfully *bool            `gorm:"default:null"`

	Message            string `gorm:"default:"""`
	CompilationMessage string

	UserId    uint
	ProblemId uint

	SubmissionTests []SubmissionTest `gorm:"foreignKey:SubmissionId;references:ID;constraint:OnDelete:CASCADE";json:"-"`
}
