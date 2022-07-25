package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/marius004/phoenix/entities"
	"github.com/marius004/phoenix/internal"
	"github.com/marius004/phoenix/models"
)

func (api *API) login(w http.ResponseWriter, r *http.Request) {
	var data models.LoginRequest

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

	user, err := api.services.UserService.GetUserByUsername(r.Context(), data.Username)
	if user == nil || err != nil {
		errorResponse(w, internal.ErrUserNotFound.Error(), http.StatusNotFound)
		return
	}

	if match := internal.CompareHashAndPassword(data.Password, user.Password); !match {
		errorResponse(w, internal.ErrInvalidPassword.Error(), http.StatusBadRequest)
		return
	}

	duration := time.Duration(api.config.CookieLifetime)
	token, err := internal.GenerateJwtToken(api.config.JwtSecret, duration, user)

	if err != nil {
		fmt.Println("Could not generate JWT token", err)
		errorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	api.setAuthCookie(w, duration, token)
	okResponse(w, user, http.StatusOK)
}

func (api *API) register(w http.ResponseWriter, r *http.Request) {
	var data models.SignupRequest

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

	user := newUser(data)
	err := api.services.UserService.CreateUser(r.Context(), user)

	if errors.Is(err, internal.ErrUsernameAlreadyExists) || errors.Is(err, internal.ErrEmailAlreadyExists) {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	} else if err != nil {
		errorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	duration := time.Duration(api.config.CookieLifetime)
	token, err := internal.GenerateJwtToken(api.config.JwtSecret, duration, user)

	if err != nil {
		fmt.Println("Could not generate JWT token", err)
		errorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	api.setAuthCookie(w, duration, token)
	okResponse(w, user, http.StatusOK)
}

func (api *API) logout(w http.ResponseWriter, r *http.Request) {
	api.deleteAuthCookie(w)
	emptyResponse(w, http.StatusOK)
}

func newUser(data models.SignupRequest) *entities.User {
	return &entities.User{
		Username: data.Username,
		Email:    data.Email,
		Password: data.Password,
	}
}
