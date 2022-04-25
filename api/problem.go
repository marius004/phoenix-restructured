package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/marius004/phoenix/internal"
	"github.com/marius004/phoenix/models"
)

func (api *API) getProblems(w http.ResponseWriter, r *http.Request) {
	filter := api.parseProblemFilter(r.URL)
	problems, err := api.services.ProblemService.GetProblemsByFilter(filter)

	if err != nil {
		fmt.Println(err)
		errorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	okResponse(w, problems, http.StatusOK)
}

func (api *API) getProblemByName(w http.ResponseWriter, r *http.Request) {
	problem := problemFromRequestContext(r.Context())

	if problem == nil {
		errorResponse(w, internal.ErrProblemDoesNotExist.Error(), http.StatusNotFound)
		return
	}

	okResponse(w, problem, http.StatusOK)
}

func (api *API) createProblem(w http.ResponseWriter, r *http.Request) {
	var data models.CreateProblemRequest

	jsonDecoder := json.NewDecoder(r.Body)
	jsonDecoder.DisallowUnknownFields()

	if err := jsonDecoder.Decode(&data); err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := data.Validate(); err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	author := userFromRequestContext(r.Context())
	problem := models.NewProblem(data, author.ID)
	err := api.services.ProblemService.CreateProblem(problem)

	if errors.Is(err, internal.ErrProblemNameAlreadyExists) {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	} else if err != nil {
		errorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	emptyResponse(w, http.StatusCreated)
}

func (api *API) updateProblemByName(w http.ResponseWriter, r *http.Request) {
	problem := problemFromRequestContext(r.Context())
	user := userFromRequestContext(r.Context())

	var data *models.UpdateProblemRequest

	jsonDecoder := json.NewDecoder(r.Body)
	jsonDecoder.DisallowUnknownFields()

	if err := jsonDecoder.Decode(&data); err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !api.canManageProblem(problem, user) {
		emptyResponse(w, http.StatusUnauthorized)
		return
	}

	if err := data.Validate(); err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := api.services.ProblemService.UpdateProblemByID(problem.ID, user, data)
	if errors.Is(err, internal.ErrUnauthorized) {
		errorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if errors.Is(err, internal.ErrProblemDoesNotExist) {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	emptyResponse(w, http.StatusOK)
}

func (api *API) deleteProblem(w http.ResponseWriter, r *http.Request) {
	problem := problemFromRequestContext(r.Context())
	user := userFromRequestContext(r.Context())

	if !api.canManageProblem(problem, user) {
		emptyResponse(w, http.StatusUnauthorized)
		return
	}

	if err := api.services.ProblemService.DeleteProblem(problem); err != nil {
		errorResponse(w, internal.ErrCouldNotDeleteProblem.Error(), http.StatusInternalServerError)
		return
	}

	emptyResponse(w, http.StatusOK)
}

func (api *API) parseProblemFilter(url *url.URL) *models.ProblemFilter {
	ret := models.ProblemFilter{}

	if v, ok := url.Query()["authorId"]; ok {
		ret.AuthorId = parseStringToUnsignedInt(v[0])
	}

	if v, ok := url.Query()["limit"]; ok {
		ret.Limit = parseStringToUnsignedInt(v[0])

		if ret.Limit > 30 {
			ret.Limit = 30
		}
	}

	return &ret
}

func parseStringToUnsignedInt(s string) uint {
	value, err := strconv.Atoi(s)

	if err != nil {
		return 0
	}

	return uint(value)
}
