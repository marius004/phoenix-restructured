package internal

import "errors"

var (
	ErrInvalidPassword       = errors.New("invalid username or password")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrUserNotFound          = errors.New("user does not exist")

	ErrProblemTestDoesNotExist    = errors.New("problem test does not exist")
	ErrSubmissionDoesNotExist     = errors.New("submission does not exist")
	ErrSubmissionTestDoesNotExist = errors.New("submission test does not exist")

	ErrEmailAlreadyExists       = errors.New("email already exists")
	ErrProblemDoesNotExist      = errors.New("problem does not exist")
	ErrProblemNameAlreadyExists = errors.New("problem name already exists")

	ErrCouldNotGeneratePasswordHash = errors.New("could not generate password hash")

	ErrUnauthorized    = errors.New("unauthorized")
	ErrMustNotBeAuthed = errors.New("you must not be authed")
	ErrMustBeAuthed    = errors.New("you must be authed")

	ErrMustBeAdmin    = errors.New("you must have admin role to do this action")
	ErrMustBeProposer = errors.New("you must have proposer role to do this action")

	ErrLangNotFound = errors.New("we currently do not provide the specified language")
)
