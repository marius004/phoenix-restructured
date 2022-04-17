package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/marius004/phoenix/entities"
)

var (
	submissionLanguageValidation   = []validation.Rule{validation.Required, validation.In(entities.CPP)}
	submissionProblemIdValidation  = []validation.Rule{validation.Required, validation.Min(0)}
	submissionSourceCodeValidation = []validation.Rule{validation.Required, validation.Length(1, 0)}
)

type CreateSubmissionRequest struct {
	Language   entities.ProgrammingLanguage
	ProblemId  int
	SourceCode []byte
}

func (data CreateSubmissionRequest) Validate() error {
	return validation.ValidateStruct(&data,
		validation.Field(&data.Language, submissionLanguageValidation...),
		validation.Field(&data.ProblemId, submissionProblemIdValidation...),
		validation.Field(&data.SourceCode, submissionSourceCodeValidation...),
	)
}

type SubmissionFilter struct {
	UserId              int
	ProblemId           int
	Score               int
	Status              entities.SubmissionStatus
	CompiledSuccesfully *bool

	Limit  int
	Offset int
}
