package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

var (
	testScoreValidation  = []validation.Rule{validation.Required, validation.Min(1), validation.Max(100)}
	testInputValidation  = []validation.Rule{validation.Required}
	testOutputValidation = []validation.Rule{validation.Required}
)

type CreateProblemTestRequest struct {
	Score  int    `json:"score"`
	Input  []byte `json:"input"`
	Output []byte `json:"output"`
}

type UpdateProblemTestRequest struct {
	Score  int    `json:"score"`
	Input  []byte `json:"input"`
	Output []byte `json:"output"`
}

func (data CreateProblemTestRequest) Validate() error {
	return validation.ValidateStruct(&data,
		validation.Field(&data.Score, testScoreValidation...),
		validation.Field(&data.Input, testInputValidation...),
		validation.Field(&data.Output, testOutputValidation...),
	)
}

func (data UpdateProblemTestRequest) Validate() error {
	return validation.ValidateStruct(&data,
		validation.Field(&data.Score, testScoreValidation...),
		validation.Field(&data.Input, testInputValidation...),
		validation.Field(&data.Output, testOutputValidation...),
	)
}
