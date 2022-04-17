package models

import "github.com/marius004/phoenix/entities"

type CreateSubmissionRequest struct {
	Lang       entities.ProgrammingLanguage
	ProblemId  uint
	SourceCode []byte
}
