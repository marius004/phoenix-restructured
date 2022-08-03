package api

import (
	"encoding/json"
	"net/http"

	"github.com/marius004/phoenix-/entities"
	"github.com/marius004/phoenix-/internal"
	"github.com/marius004/phoenix-/models"
)

func (api *API) getProblemTests(w http.ResponseWriter, r *http.Request) {
	problem := internal.ProblemFromContext(r.Context())
	tests, err := api.services.ProblemTestService.GetProblemTestsByProblemID(r.Context(), problem.ID)

	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	okResponse(w, tests, http.StatusOK)
}

func (api *API) getProblemTestByID(w http.ResponseWriter, r *http.Request) {
	user := internal.UserFromContext(r.Context())
	problem := internal.ProblemFromContext(r.Context())

	if !internal.CanManageProblem(problem, user) {
		errorResponse(w, internal.ErrUnauthorized.Error(), http.StatusUnauthorized)
		return
	}

	problemTest := internal.ProblemTestFromContext(r.Context())
	okResponse(w, problemTest, http.StatusOK)
}

func (api *API) updateProblemTestById(w http.ResponseWriter, r *http.Request) {
	user := internal.UserFromContext(r.Context())

	problem := internal.ProblemFromContext(r.Context())
	problemTest := internal.ProblemTestFromContext(r.Context())

	if !internal.CanManageProblem(problem, user) {
		errorResponse(w, internal.ErrUnauthorized.Error(), http.StatusUnauthorized)
		return
	}

	jsonDecoder := json.NewDecoder(r.Body)
	jsonDecoder.DisallowUnknownFields()

	var data models.UpdateProblemTestRequest
	if err := jsonDecoder.Decode(&data); err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := data.Validate(); err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := api.services.ProblemTestService.UpdateProblemTestByID(r.Context(), problemTest.ID, &data); err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	emptyResponse(w, http.StatusOK)
}

func (api *API) deleteProblemTestById(w http.ResponseWriter, r *http.Request) {
	user := internal.UserFromContext(r.Context())

	problem := internal.ProblemFromContext(r.Context())
	problemTest := internal.ProblemTestFromContext(r.Context())

	if !internal.CanManageProblem(problem, user) {
		errorResponse(w, internal.ErrUnauthorized.Error(), http.StatusUnauthorized)
		return
	}

	if err := api.services.ProblemTestService.DeleteProblemTestByID(r.Context(), problemTest.ID); err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	emptyResponse(w, http.StatusOK)
}

func (api *API) createProblemTest(w http.ResponseWriter, r *http.Request) {
	problem := internal.ProblemFromContext(r.Context())
	user := internal.UserFromContext(r.Context())

	jsonDecoder := json.NewDecoder(r.Body)
	jsonDecoder.DisallowUnknownFields()

	var data models.CreateProblemTestRequest
	if err := jsonDecoder.Decode(&data); err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !internal.CanManageProblem(problem, user) {
		errorResponse(w, internal.ErrUnauthorized.Error(), http.StatusUnauthorized)
		return
	}

	if err := data.Validate(); err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	problemTest := entities.ProblemTest{}
	problemTest.ProblemId = problem.ID

	problemTest.Score = data.Score
	problemTest.Input = data.Input
	problemTest.Output = data.Output

	if err := api.services.ProblemTestService.CreateProblemTest(r.Context(), &problemTest); err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	emptyResponse(w, http.StatusCreated)
}
