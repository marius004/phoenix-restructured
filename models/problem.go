package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/marius004/phoenix/entities"
)

const (
	Easy    = "easy"
	Medium  = "medium"
	Hard    = "hard"
	Contest = "contest"
)

var (
	problemNameValidation        = []validation.Rule{validation.Required, validation.Length(3, 20)}
	problemTimeLimitValidation   = []validation.Rule{validation.Required, validation.Min(0.0), validation.Max(2.0)}
	problemMemoryLimitValidation = []validation.Rule{validation.Required, validation.Min(0), validation.Max(65537)}
	problemStackLimitValidation  = []validation.Rule{validation.Required, validation.Min(0), validation.Max(16384)}
	problemDifficultyValidation  = []validation.Rule{validation.Required, validation.In(Easy, Medium, Hard, Contest)}
)

type ProblemFilter struct {
	AuthorId uint
	Limit    uint
}

type CreateProblemRequest struct {
	Name        string
	Description string
	Difficulty  string

	TimeLimit   float32
	MemoryLimit int
	StackLimit  int
}

type UpdateProblemRequest struct {
	Description string
	Difficulty  string

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
		validation.Field(&data.Difficulty, problemDifficultyValidation...),
	)
}

func (data UpdateProblemRequest) Validate() error {
	return validation.ValidateStruct(&data,
		validation.Field(&data.TimeLimit, problemTimeLimitValidation...),
		validation.Field(&data.MemoryLimit, problemMemoryLimitValidation...),
		validation.Field(&data.StackLimit, problemStackLimitValidation...),
		validation.Field(&data.Difficulty, problemDifficultyValidation...),
	)
}

func NewProblem(request CreateProblemRequest, authorId uint) *entities.Problem {
	return &entities.Problem{
		Name:        request.Name,
		Description: request.Description,
		AuthorId:    authorId,
		Difficulty:  request.Difficulty,

		TimeLimit:   request.TimeLimit,
		MemoryLimit: request.MemoryLimit,
		StackLimit:  request.StackLimit,
	}
}
