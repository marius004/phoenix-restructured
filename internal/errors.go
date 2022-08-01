package internal

import "errors"

var (
	ErrCouldNotDeleteProblem = errors.New("could not delete problem")
	ErrCouldNotDeletePost    = errors.New("could not delete post")

	ErrCouldNotUpdateUser = errors.New("could not update user data")
	ErrCouldNotUpdatePost = errors.New("could not update post")
	ErrCouldNotAssignRole = errors.New("could not assign role")

	ErrInvalidPassword       = errors.New("invalid username or password")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrUserNotFound          = errors.New("user does not exist")
	ErrCouldNotGetUserStats  = errors.New("could not get user stats")

	ErrProblemTestDoesNotExist    = errors.New("problem test does not exist")
	ErrSubmissionDoesNotExist     = errors.New("submission does not exist")
	ErrSubmissionTestDoesNotExist = errors.New("submission test does not exist")
	ErrPostDoesNotExist           = errors.New("post does not exist")

	ErrEmailAlreadyExists       = errors.New("email already exists")
	ErrProblemDoesNotExist      = errors.New("problem does not exist")
	ErrProblemNameAlreadyExists = errors.New("problem name already exists")
	ErrPostTitleAlreadyExists   = errors.New("post title already exists")

	ErrCouldNotGeneratePasswordHash = errors.New("could not generate password hash")

	ErrUnauthorized    = errors.New("unauthorized")
	ErrMustNotBeAuthed = errors.New("you must not be authed")
	ErrMustBeAuthed    = errors.New("you must be authed")

	ErrMustBeAdmin    = errors.New("you must have admin role to do this action")
	ErrMustBeProposer = errors.New("you must have proposer role to do this action")

	ErrLangNotFound = errors.New("we currently do not provide the specified language")
)
