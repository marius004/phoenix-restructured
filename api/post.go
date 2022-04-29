package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/marius004/phoenix/entities"
	"github.com/marius004/phoenix/internal"
	"github.com/marius004/phoenix/models"
)

func (api *API) getPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := api.services.PostService.GetPosts()

	if err != nil {
		errorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	okResponse(w, posts, http.StatusOK)
}

func (api *API) getPost(w http.ResponseWriter, r *http.Request) {
	post := postFromRequestContext(r.Context())
	okResponse(w, post, http.StatusOK)
}

func (api *API) createPost(w http.ResponseWriter, r *http.Request) {
	author := userFromRequestContext(r.Context())
	var data models.CreatePostRequest

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

	post := &entities.Post{
		Title:   data.Title,
		Content: []byte(data.Content),

		UserId: author.ID,
	}

	err := api.services.PostService.CreatePost(post)
	if errors.Is(err, internal.ErrPostTitleAlreadyExists) {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err != nil {
		errorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	emptyResponse(w, http.StatusOK)
}

func (api *API) updatePost(w http.ResponseWriter, r *http.Request) {
	var data models.UpdatePostRequest

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
	post := postFromRequestContext(r.Context())

	if !api.canManagePost(post, author) {
		errorResponse(w, internal.ErrUnauthorized.Error(), http.StatusUnauthorized)
		return
	}

	if err := api.services.PostService.UpdatePostByID(post.ID, &data); err != nil {
		errorResponse(w, internal.ErrCouldNotUpdatePost.Error(), http.StatusInternalServerError)
		return
	}

	emptyResponse(w, http.StatusOK)
}

func (api *API) deletePost(w http.ResponseWriter, r *http.Request) {
	author := userFromRequestContext(r.Context())
	post := postFromRequestContext(r.Context())

	if !api.canManagePost(post, author) {
		errorResponse(w, internal.ErrUnauthorized.Error(), http.StatusUnauthorized)
		return
	}

	if err := api.services.PostService.DeletePost(post); err != nil {
		errorResponse(w, internal.ErrCouldNotDeletePost.Error(), http.StatusInternalServerError)
		return
	}

	emptyResponse(w, http.StatusOK)
}
