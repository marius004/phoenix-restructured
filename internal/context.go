package internal

import (
	"context"

	"github.com/marius004/phoenix-/entities"
)

const (
	UserContextKey        = "user"
	ProblemContextKey     = "problem"
	ProblemTestContextKey = "problemTest"
	SubmissionContextKey  = "submission"
)

func UserFromContext(context context.Context) *entities.User {
	switch usr := context.Value(UserContextKey).(type) {
	case entities.User:
		return &usr
	case *entities.User:
		return usr
	default:
		return nil
	}
}

func ProblemFromContext(context context.Context) *entities.Problem {
	switch problem := context.Value(ProblemContextKey).(type) {
	case entities.Problem:
		return &problem
	case *entities.Problem:
		return problem
	default:
		return nil
	}
}

func ProblemTestFromContext(context context.Context) *entities.ProblemTest {
	switch problemTest := context.Value(ProblemTestContextKey).(type) {
	case entities.ProblemTest:
		return &problemTest
	case *entities.ProblemTest:
		return problemTest
	default:
		return nil
	}
}

func SubmissionFromContext(context context.Context) *entities.Submission {
	switch submission := context.Value(SubmissionContextKey).(type) {
	case entities.Submission:
		return &submission
	case *entities.Submission:
		return submission
	default:
		return nil
	}
}
