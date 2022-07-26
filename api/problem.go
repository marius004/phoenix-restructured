package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/marius004/phoenix/entities"
	"github.com/marius004/phoenix/internal"
	"github.com/marius004/phoenix/models"
)

func (api *API) getProblems(w http.ResponseWriter, r *http.Request) {
	filter := api.parseProblemFilter(r)

	problems, err := api.services.ProblemService.GetProblemsByFilter(r.Context(), filter)

	if err != nil {
		fmt.Println(err)
		errorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	okResponse(w, problems, http.StatusOK)
}

func (api *API) getProblemByName(w http.ResponseWriter, r *http.Request) {
	problem := internal.ProblemFromContext(r.Context())

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

	author := internal.UserFromContext(r.Context())
	problem := models.NewProblem(data, author.ID)
	err := api.services.ProblemService.CreateProblem(r.Context(), problem)

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
	problem := internal.ProblemFromContext(r.Context())
	user := internal.UserFromContext(r.Context())

	var data *models.UpdateProblemRequest

	jsonDecoder := json.NewDecoder(r.Body)
	jsonDecoder.DisallowUnknownFields()

	if err := jsonDecoder.Decode(&data); err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !internal.CanManageProblem(problem, user) {
		emptyResponse(w, http.StatusUnauthorized)
		return
	}

	if err := data.Validate(); err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := api.services.ProblemService.UpdateProblemByID(r.Context(), problem.ID, user, data)
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
	problem := internal.ProblemFromContext(r.Context())
	user := internal.UserFromContext(r.Context())

	if !internal.CanManageProblem(problem, user) {
		emptyResponse(w, http.StatusUnauthorized)
		return
	}

	if err := api.services.ProblemService.DeleteProblem(r.Context(), problem); err != nil {
		errorResponse(w, internal.ErrCouldNotDeleteProblem.Error(), http.StatusInternalServerError)
		return
	}

	emptyResponse(w, http.StatusOK)
}

func (api *API) publishProblem(w http.ResponseWriter, r *http.Request) {
	user := internal.UserFromContext(r.Context())
	problem := internal.ProblemFromContext(r.Context())

	if !internal.CanManageProblem(problem, user) {
		errorResponse(w, internal.ErrUnauthorized.Error(), http.StatusUnauthorized)
		return
	}

	var status entities.ProblemStatus
	if internal.IsUserAdmin(user) {
		status = entities.Published
	} else { // proposer
		status = entities.WaitingForApproval
	}

	if err := api.services.ProblemService.UpdateProblemStatus(r.Context(), problem, status); err != nil {
		errorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if status == entities.Published {
		okResponse(w, models.NewUpdateProblemStatusResponse("Problem published succesfully", status), http.StatusOK)
	} else {
		okResponse(w, models.NewUpdateProblemStatusResponse("Problem waiting for admin approval", status), http.StatusOK)
	}
}

func (api *API) unpublishProblem(w http.ResponseWriter, r *http.Request) {
	user := internal.UserFromContext(r.Context())
	problem := internal.ProblemFromContext(r.Context())

	if !internal.CanManageProblem(problem, user) {
		errorResponse(w, internal.ErrUnauthorized.Error(), http.StatusUnauthorized)
		return
	}

	if err := api.services.ProblemService.UpdateProblemStatus(r.Context(), problem, entities.UnPublished); err != nil {
		errorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	okResponse(w, models.NewUpdateProblemStatusResponse("Problem unpublished succesfully", entities.UnPublished), http.StatusOK)
}

func (api *API) parseProblemFilter(r *http.Request) *models.ProblemFilter {
	url := r.URL
	ret := models.ProblemFilter{}

	if v, ok := url.Query()["authorId"]; ok {
		ret.AuthorId = parseStringToUnsignedInt(v[0])
	}

	if v, ok := url.Query()["problemId"]; ok {
		ret.ProblemId = parseStringToUnsignedInt(v[0])
	}

	if v, ok := url.Query()["limit"]; ok {
		ret.Limit = int(parseStringToUnsignedInt(v[0]))

		if ret.Limit > 30 {
			ret.Limit = 30
		}
	}

	if v, ok := url.Query()["status"]; ok {
		if status := entities.ProblemStatus(v[0]); status.IsValid() {
			ret.Status = status
		}
	}

	// overwrite the filter values
	if user := internal.UserFromContext(r.Context()); !internal.IsUserProposer(user) {
		ret.Status = entities.Published
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
