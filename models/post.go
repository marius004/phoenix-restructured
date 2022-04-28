package models

import validation "github.com/go-ozzo/ozzo-validation"

var (
	titleValidation   = []validation.Rule{validation.Required, validation.Length(1, 50)}
	contentValidation = []validation.Rule{validation.Required, validation.Length(1, 0)}
)

type CreatePostRequest struct {
	Title   string
	Content string
}

type UpdatePostRequest struct {
	Title   string
	Content string
}

func (data CreatePostRequest) Validate() error {
	return validation.ValidateStruct(&data,
		validation.Field(&data.Title, titleValidation...),
		validation.Field(&data.Content, contentValidation...),
	)
}

func (data UpdatePostRequest) Validate() error {
	return validation.ValidateStruct(&data,
		validation.Field(&data.Title, titleValidation...),
		validation.Field(&data.Content, contentValidation...),
	)
}
