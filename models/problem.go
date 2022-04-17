package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/marius004/phoenix/entities"
)

var (
	problemNameValidation        = []validation.Rule{validation.Required, validation.Length(3, 20)}
	problemTimeLimitValidation   = []validation.Rule{validation.Required, validation.Min(0.0), validation.Max(2.0)}
	problemMemoryLimitValidation = []validation.Rule{validation.Required, validation.Min(0), validation.Max(65537)}
	problemStackLimitValidation  = []validation.Rule{validation.Required, validation.Min(0), validation.Max(16384)}
)

type ProblemFilter struct {
	AuthorId uint
	Limit    uint
}

type CreateProblemRequest struct {
	Name        string
	Description string

	TimeLimit   float32
	MemoryLimit int
	StackLimit  int
}

type UpdateProblemRequest struct {
	Description string

	TimeLimit   float32
	MemoryLimit int
	StackLimit  int
}

func (data CreateProblemRequest) Validate() error {
	return validation.ValidateStruct(&data,
		validation.Field(&data.Name, problemNameValidation...),
		validation.Field(&data.TimeLimit, problemTimeLimitValidation...),
		validation.Field(&data.MemoryLimit, problemMemoryLimitValidation...),
		validation.Field(&data.StackLimit, problemStackLimitValidation...),
	)
}

func NewProblem(request CreateProblemRequest, authorId uint) *entities.Problem {
	return &entities.Problem{
		Name:        request.Name,
		Description: request.Description,
		AuthorId:    authorId,

		TimeLimit:   request.TimeLimit,
		MemoryLimit: request.MemoryLimit,
		StackLimit:  request.StackLimit,
	}
}
