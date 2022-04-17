package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/marius004/phoenix/entities"
	"github.com/marius004/phoenix/internal"
	"github.com/marius004/phoenix/models"
)

func (api *API) getSubmissions(w http.ResponseWriter, r *http.Request) {
	filter := api.parseSubmissionFilter(r)
	submissions, err := api.services.SubmissionService.GetBySubmissionFilter(*filter)

	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	okResponse(w, submissions, http.StatusOK)
}

func (api *API) getSubmissionById(w http.ResponseWriter, r *http.Request) {
	submission := submissionFromRequestContext(r.Context())
	okResponse(w, submission, http.StatusOK)
}

func (api *API) createSubmission(w http.ResponseWriter, r *http.Request) {
	user := userFromRequestContext(r.Context())
	var data models.CreateSubmissionRequest

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

	if problem, err := api.services.ProblemService.GetProblemByID(uint(data.ProblemId)); err != nil || problem == nil {
		errorResponse(w, internal.ErrProblemDoesNotExist.Error(), http.StatusBadRequest)
		return
	}

	submission := &entities.Submission{}

	submission.ProblemId = uint(data.ProblemId)
	submission.UserId = user.ID

	submission.SourceCode = data.SourceCode
	submission.Language = data.Language
	submission.Status = entities.Waiting

	if err := api.services.SubmissionService.CreateSubmission(submission); err != nil {
		errorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	emptyResponse(w, http.StatusCreated)
}

func (api *API) parseSubmissionFilter(r *http.Request) *models.SubmissionFilter {
	ret := models.SubmissionFilter{}

	if v, ok := r.URL.Query()["username"]; ok {
		username := v[0]

		user, err := api.services.UserService.GetUserByUsername(username)
		if err == nil && user != nil {
			ret.UserId = int(user.ID)
		} else {
			ret.UserId = -1
		}
	}

	if v, ok := r.URL.Query()["userId"]; ok {
		userId, err := strconv.Atoi(v[0])

		if err == nil && userId > 0 {
			ret.UserId = userId
		} else {
			ret.UserId = -1
		}
	}

	if v, ok := r.URL.Query()["problem"]; ok {
		problemName := v[0]

		problem, err := api.services.ProblemService.GetProblemByName(problemName)
		if err == nil && problem != nil {
			ret.ProblemId = int(problem.ID)
		} else {
			ret.ProblemId = -1
		}
	}

	if v, ok := r.URL.Query()["problemId"]; ok {
		problemId, err := strconv.Atoi(v[0])

		if err == nil && problemId > 0 {
			ret.ProblemId = problemId
		} else {
			ret.ProblemId = -1
		}
	}

	if v, ok := r.URL.Query()["score"]; ok {
		score, err := strconv.Atoi(v[0])

		if err == nil && score > 0 {
			ret.Score = score
		} else {
			ret.Score = -1
		}
	}

	if v, ok := r.URL.Query()["status"]; ok {
		ret.Status = entities.SubmissionStatus(v[0])
	}

	if v, ok := r.URL.Query()["isCompiled"]; ok {
		compiled, err := strconv.ParseBool(v[0])

		if err != nil {
			ret.CompiledSuccesfully = &compiled
		} else {
			ret.CompiledSuccesfully = nil
		}
	}

	if v, ok := r.URL.Query()["limit"]; ok {
		if limit, err := strconv.Atoi(v[0]); err != nil {
			ret.Limit = limit
		} else {
			ret.Limit = -1
		}
	}

	if v, ok := r.URL.Query()["offset"]; ok {
		if offset, err := strconv.Atoi(v[0]); err != nil {
			ret.Offset = offset
		} else {
			ret.Offset = -1
		}
	}

	return &ret
}
