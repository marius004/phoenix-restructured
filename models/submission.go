package models

import "github.com/marius004/phoenix/entities"

type CreateSubmissionRequest struct {
	Lang       entities.ProgrammingLanguage
	ProblemId  uint
	SourceCode []byte
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
