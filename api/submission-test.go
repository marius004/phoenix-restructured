package api

import "net/http"

func (api *API) getSubmissionTests(w http.ResponseWriter, r *http.Request) {
	submission := submissionFromRequestContext(r.Context())
	submissionTests, err := api.services.SubmissionTestService.GetSubmissionTestsBySubmissionID(r.Context(), submission.ID)

	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	okResponse(w, submissionTests, http.StatusOK)
}
