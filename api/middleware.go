package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/v5"
	"github.com/marius004/phoenix-algo/internal"
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

			user, err := api.services.UserService.GetUserByID(r.Context(), uint(userId))

			if err != nil {
				fmt.Println(err)
				next.ServeHTTP(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), internal.UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))

			return
		}

		next.ServeHTTP(w, r)
	})
}

func (api *API) problemCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "problemName")
		problem, err := api.services.ProblemService.GetProblemByName(r.Context(), name)

		if err != nil || problem == nil {
			errorResponse(w, internal.ErrProblemDoesNotExist.Error(), http.StatusNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), internal.ProblemContextKey, problem)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (api *API) problemTestCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, _ := internal.ConvertStringToUint(chi.URLParam(r, "problemTestId"))
		problemTest, err := api.services.ProblemTestService.GetProblemTestByID(r.Context(), id)

		if err != nil || problemTest == nil {
			errorResponse(w, internal.ErrProblemTestDoesNotExist.Error(), http.StatusNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), internal.ProblemTestContextKey, problemTest)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (api *API) submissionCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, _ := internal.ConvertStringToUint(chi.URLParam(r, "submissionId"))
		submission, err := api.services.SubmissionService.GetSubmissionByID(r.Context(), id)

		if err != nil || submission == nil {
			errorResponse(w, internal.ErrSubmissionDoesNotExist.Error(), http.StatusNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), internal.SubmissionContextKey, submission)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (api *API) mustBeAuthed(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user := internal.UserFromContext(r.Context()); !internal.IsUserAuthed(user) {
			errorResponse(w, internal.ErrUnauthorized.Error(), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (api *API) mustNotBeAuthed(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user := internal.UserFromContext(r.Context()); internal.IsUserAuthed(user) {
			errorResponse(w, internal.ErrMustNotBeAuthed.Error(), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (api *API) mustBeAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user := internal.UserFromContext(r.Context()); !internal.IsUserAdmin(user) {
			errorResponse(w, internal.ErrMustBeAdmin.Error(), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (api *API) mustBeProposer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user := internal.UserFromContext(r.Context()); !internal.IsUserProposer(user) {
			errorResponse(w, internal.ErrMustBeProposer.Error(), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
