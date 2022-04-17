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

	r.Route("/auth", func(r chi.Router) {
		r.With(api.mustNotBeAuthed).Post("/register", api.register)
		r.With(api.mustNotBeAuthed).Post("/login", api.login)
		r.With(api.mustBeAuthed).Post("/logout", api.logout)
	})

	r.Route("/problems", func(r chi.Router) {
		r.Get("/", api.getProblems)
		r.With(api.mustBeProposer).Post("/", api.createProblem)

		r.With(api.problemCtx).Route("/{problemId}", func(r chi.Router) {
			r.Get("/", api.getProblemByID)
			r.With(api.mustBeProposer).Put("/", api.updateProblemByID)

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

	return r
}

func NewAPI(config *internal.Config, services *internal.Services) *API {
	return &API{
		config:   config,
		services: services,
	}
}
