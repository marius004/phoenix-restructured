package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/marius004/phoenix/internal"
	"github.com/marius004/phoenix/models"
)

func (api *API) getUserByUsername(w http.ResponseWriter, r *http.Request) {
	logger := internal.GetGlobalLoggerInstance()
	username := chi.URLParam(r, "username")

	user, err := api.services.UserService.GetUserByUsername(r.Context(), username)

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

func (api *API) getUsers(w http.ResponseWriter, r *http.Request) {
	filter := api.parseUserFilter(r)
	users, err := api.services.UserService.GetUsers(r.Context(), filter)

	if err != nil {
		errorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	okResponse(w, users, http.StatusOK)
}

func (api *API) updateUser(w http.ResponseWriter, r *http.Request) {
	user := internal.UserFromContext(r.Context())

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

	if err := api.services.UserService.UpdateUser(r.Context(), user, &data); err != nil {
		errorResponse(w, internal.ErrCouldNotUpdateUser.Error(), http.StatusInternalServerError)
		return
	}

	emptyResponse(w, http.StatusOK)
}

func (api *API) assignProposerRole(w http.ResponseWriter, r *http.Request) {
	value, err := strconv.ParseBool(chi.URLParam(r, "value"))
	username := chi.URLParam(r, "username")

	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = api.services.UserService.AssignProposerRole(r.Context(), username, value)
	if err != nil {
		errorResponse(w, internal.ErrCouldNotAssignRole.Error(), http.StatusInternalServerError)
		return
	}

	emptyResponse(w, http.StatusOK)
}

func (api *API) parseUserFilter(r *http.Request) *models.UserFilter {
	ret := models.UserFilter{}

	if v, ok := r.URL.Query()["userId"]; ok {
		userId, err := strconv.Atoi(v[0])

		if err == nil && userId > 0 {
			ret.UserId = userId
		} else {
			ret.UserId = -1
		}
	}

	if v, ok := r.URL.Query()["username"]; ok {
		ret.Username = v[0]
	}

	if v, ok := r.URL.Query()["email"]; ok {
		ret.Email = v[0]
	}

	if v, ok := r.URL.Query()["linkedInUrl"]; ok {
		ret.LinkedInURL = v[0]
	}

	if v, ok := r.URL.Query()["githubUrl"]; ok {
		ret.GithubURL = v[0]
	}

	if v, ok := r.URL.Query()["websiteUrl"]; ok {
		ret.WebsiteURL = v[0]
	}

	if v, ok := r.URL.Query()["isAdmin"]; ok {
		value, err := strconv.ParseBool(v[0])

		if err == nil {
			ret.IsAdmin = &value
		}
	}

	if v, ok := r.URL.Query()["isProposer"]; ok {
		value, err := strconv.ParseBool(v[0])

		if err == nil {
			ret.IsProposer = &value
		}
	}

	if v, ok := r.URL.Query()["userIconUrl"]; ok {
		ret.UserIconURL = v[0]
	}

	return &ret
}
