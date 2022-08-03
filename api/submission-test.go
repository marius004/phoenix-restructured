package api

import (
	"net/http"

	"github.com/marius004/phoenix-algo/internal"
)

func (api *API) getSubmissionTests(w http.ResponseWriter, r *http.Request) {
	submission := internal.SubmissionFromContext(r.Context())
	submissionTests, err := api.services.SubmissionTestService.GetSubmissionTestsBySubmissionID(r.Context(), submission.ID)

	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	okResponse(w, submissionTests, http.StatusOK)
}
