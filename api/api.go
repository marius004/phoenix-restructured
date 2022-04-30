package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/marius004/phoenix/internal"
)

type API struct {
	config   *internal.Config
	services *internal.Services
}

func (api *API) Routes() http.Handler {
	r := chi.NewRouter()

	r.Use(api.jwtMiddleware)
	go api.services.Grader.Handle()

	r.Route("/auth", func(r chi.Router) {
		r.With(api.mustNotBeAuthed).Post("/register", api.register)
		r.With(api.mustNotBeAuthed).Post("/login", api.login)
		r.With(api.mustBeAuthed).Post("/logout", api.logout)
	})

	r.Route("/users", func(r chi.Router) {
		r.Get("/", api.getUsers)
		r.Get("/{username}", api.getUserByUsername)
		r.With(api.mustBeAuthed).Put("/", api.updateUser)
	})

	r.Route("/problems", func(r chi.Router) {
		r.Get("/", api.getProblems)
		r.With(api.mustBeProposer).Post("/", api.createProblem)

		r.With(api.problemCtx).Route("/{problemName}", func(r chi.Router) {
			r.Get("/", api.getProblemByName)
			r.With(api.mustBeProposer).Put("/", api.updateProblemByName)
			r.With(api.mustBeProposer).Delete("/", api.deleteProblem)

			r.Route("/tests", func(r chi.Router) {
				r.With(api.mustBeProposer).Get("/", api.getProblemTests)
				r.With(api.mustBeProposer).Post("/", api.createProblemTest)

				r.With(api.problemTestCtx).Route("/{problemTestId}", func(r chi.Router) {
					r.With(api.mustBeProposer).Get("/", api.getProblemTestByID)
					r.With(api.mustBeProposer).Put("/", api.updateProblemTestById)
					r.With(api.mustBeProposer).Delete("/", api.deleteProblemTestById)
				})
			})
		})
	})

	r.Route("/submissions", func(r chi.Router) {
		r.Get("/", api.getSubmissions)
		r.With(api.mustBeAuthed).Post("/", api.createSubmission)

		r.With(api.submissionCtx).Route("/{submissionId}", func(r chi.Router) {
			r.Get("/", api.getSubmissionById)
			r.Get("/tests", api.getSubmissionTests)
		})
	})

	r.Route("/posts", func(r chi.Router) {
		r.Get("/", api.getPosts)
		r.With(api.mustBeProposer).Post("/", api.createPost)

		r.With(api.postCtx).Route("/{postId}", func(r chi.Router) {
			r.Get("/", api.getPost)
			r.With(api.mustBeProposer).Put("/", api.updatePost)
			r.With(api.mustBeProposer).Delete("/", api.deletePost)
		})
	})

	return r
}

func NewAPI(config *internal.Config, services *internal.Services) *API {
	return &API{
		config:   config,
		services: services,
	}
}
