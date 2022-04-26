package models

import validation "github.com/go-ozzo/ozzo-validation"

type UpdateUserRequest struct {
	Bio string

	LinkedInURL string `json:"linkedinURL"`
	GithubURL   string `json:"githubURL"`
	WebsiteURL  string `json:"websiteURL"`

	UserIconURL string `json:"userIconURL"`
}

var bioValidation = []validation.Rule{validation.Length(0, 250)}

func (data UpdateUserRequest) Validate() error {
	return validation.ValidateStruct(&data,
		validation.Field(&data.Bio, bioValidation...),
	)
}

type UserFilter struct {
	UserId   int
	Username string
	Email    string

	LinkedInURL string
	GithubURL   string
	WebsiteURL  string
	UserIconURL string

	IsAdmin    *bool
	IsProposer *bool

	Limit  int
	Offset int
}
