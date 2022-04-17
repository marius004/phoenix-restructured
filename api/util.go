package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/marius004/phoenix/entities"
	"github.com/marius004/phoenix/internal"
)

const (
	authCookieName = "auth-token"

	userContextKey        = "user"
	problemContextKey     = "problem"
	problemTestContextKey = "problemTest"
	submissionContextKey  = "submission"
)

func errorResponse(w http.ResponseWriter, err string, status int) {
	http.Error(w, err, status)
}

func okResponse(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(data)
}

func emptyResponse(w http.ResponseWriter, status int) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
}

func (s *API) setAuthCookie(w http.ResponseWriter, duration time.Duration, token string) {
	cookie := &http.Cookie{
		Name:  authCookieName,
		Value: token,
		Path:  "/",
	}

	http.SetCookie(w, cookie)
}

func (api *API) deleteAuthCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:    authCookieName,
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
	}

	http.SetCookie(w, cookie)
}

func userFromRequestContext(context context.Context) *entities.User {
	switch usr := context.Value(userContextKey).(type) {
	case entities.User:
		return &usr
	case *entities.User:
		return usr
	default:
		return nil
	}
}

func problemFromRequestContext(context context.Context) *entities.Problem {
	switch problem := context.Value(problemContextKey).(type) {
	case entities.Problem:
		return &problem
	case *entities.Problem:
		return problem
	default:
		return nil
	}
}

func problemTestFromRequestContext(context context.Context) *entities.ProblemTest {
	switch problemTest := context.Value(problemTestContextKey).(type) {
	case entities.ProblemTest:
		return &problemTest
	case *entities.ProblemTest:
		return problemTest
	default:
		return nil
	}
}

func submissionFromRequestContext(context context.Context) *entities.Submission {
	switch submission := context.Value(submissionContextKey).(type) {
	case entities.Submission:
		return &submission
	case *entities.Submission:
		return submission
	default:
		return nil
	}
}

func convertStringToUint(s string) (uint, error) {
	res, err := strconv.Atoi(s)

	if err != nil {
		return 0, err
	}

	return uint(res), nil
}

func (api *API) canManageProblem(problem *entities.Problem, user *entities.User) bool {
	if problem == nil {
		return false
	}

	if (internal.IsUserProposer(user) && problem.AuthorId == user.ID) || internal.IsUserAdmin(user) {
		return true
	}

	return false
}
