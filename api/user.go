package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/marius004/phoenix/internal"
	"github.com/marius004/phoenix/models"
)

func (api *API) getUserByUsername(w http.ResponseWriter, r *http.Request) {
	logger := internal.GetGlobalLoggerInstance()
	username := chi.URLParam(r, "username")

	user, err := api.services.UserService.GetUserByUsername(username)

	// TODO implement unknown errror & return 500 http codes for all handlers
	if err != nil {
		logger.Println(err)
		errorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user == nil {
		emptyResponse(w, http.StatusNotFound)
		return
	}

	okResponse(w, user, http.StatusOK)
}

func (api *API) updateUser(w http.ResponseWriter, r *http.Request) {
	user := userFromRequestContext(r.Context())

	var data models.UpdateUserRequest

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

	if err := api.services.UserService.UpdateUser(user, &data); err != nil {
		errorResponse(w, internal.ErrCouldNotUpdateUser.Error(), http.StatusInternalServerError)
		return
	}

	emptyResponse(w, http.StatusOK)
}
