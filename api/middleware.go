package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/v5"
	"github.com/marius004/phoenix/internal"
)

func (api *API) jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authCookie, err := r.Cookie(authCookieName)

		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		token, err := internal.VerifyToken(authCookie.Value, api.config.JwtSecret)

		if err != nil {
			fmt.Println(err)
			next.ServeHTTP(w, r)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userId, err := strconv.Atoi(claims["iss"].(string))

			if err != nil {
				fmt.Println(err)
				next.ServeHTTP(w, r)
				return
			}

			user, err := api.services.UserService.GetUserByID(uint(userId))

			if err != nil {
				fmt.Println(err)
				next.ServeHTTP(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), userContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))

			return
		}

		next.ServeHTTP(w, r)
	})
}

func (api *API) problemCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "problemName")
		problem, err := api.services.ProblemService.GetProblemByName(name)

		if err != nil || problem == nil {
			errorResponse(w, internal.ErrProblemDoesNotExist.Error(), http.StatusNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), problemContextKey, problem)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (api *API) problemTestCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, _ := convertStringToUint(chi.URLParam(r, "problemTestId"))
		problemTest, err := api.services.ProblemTestService.GetProblemTestByID(id)

		if err != nil || problemTest == nil {
			errorResponse(w, internal.ErrProblemTestDoesNotExist.Error(), http.StatusNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), problemTestContextKey, problemTest)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (api *API) submissionCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, _ := convertStringToUint(chi.URLParam(r, "submissionId"))
		submission, err := api.services.SubmissionService.GetSubmissionByID(id)

		if err != nil || submission == nil {
			errorResponse(w, internal.ErrSubmissionDoesNotExist.Error(), http.StatusNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), submissionContextKey, submission)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (api *API) postCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postId, _ := convertStringToUint(chi.URLParam(r, "postId"))
		post, err := api.services.PostService.GetPostByID(postId)

		if err != nil || post == nil {
			errorResponse(w, internal.ErrPostDoesNotExist.Error(), http.StatusNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), postContextKey, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (api *API) testCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postId, _ := convertStringToUint(chi.URLParam(r, "testId"))
		postName := chi.URLParam(r, "testName")

		fmt.Println(postId, postName)

		next.ServeHTTP(w, r)
	})
}

func (api *API) mustBeAuthed(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user := userFromRequestContext(r.Context()); !internal.IsUserAuthed(user) {
			errorResponse(w, internal.ErrUnauthorized.Error(), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (api *API) mustNotBeAuthed(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user := userFromRequestContext(r.Context()); internal.IsUserAuthed(user) {
			errorResponse(w, internal.ErrMustNotBeAuthed.Error(), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (api *API) mustBeAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user := userFromRequestContext(r.Context()); !internal.IsUserAdmin(user) {
			errorResponse(w, internal.ErrMustBeAdmin.Error(), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (api *API) mustBeProposer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user := userFromRequestContext(r.Context()); !internal.IsUserProposer(user) {
			errorResponse(w, internal.ErrMustBeProposer.Error(), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
