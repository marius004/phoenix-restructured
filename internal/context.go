package internal

import (
	"context"
	"strconv"

	"github.com/marius004/phoenix/entities"
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

func ConvertStringToUint(s string) (uint, error) {
	res, err := strconv.Atoi(s)

	if err != nil {
		return 0, err
	}

	return uint(res), nil
}

func CanManageProblem(problem *entities.Problem, user *entities.User) bool {
	if problem == nil {
		return false
	}

	if (IsUserProposer(user) && problem.AuthorId == user.ID) || IsUserAdmin(user) {
		return true
	}

	return false
}
