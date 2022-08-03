package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/marius004/phoenix-/internal"
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
		r.With(api.mustBeAuthed).Put("/", api.updateUser)

		r.Route("/{username}", func(r chi.Router) {
			r.Get("/", api.getUserByUsername)
			r.Get("/stats", api.getUserStats)
			r.With(api.mustBeAdmin).Post("/roles/proposer/{value}", api.assignProposerRole)
		})
	})

	r.Route("/problems", func(r chi.Router) {
		r.Get("/", api.getProblems)
		r.With(api.mustBeProposer).Post("/", api.createProblem)

		r.With(api.problemCtx).Route("/{problemName}", func(r chi.Router) {
			r.Get("/", api.getProblemByName)

			r.With(api.mustBeProposer).Route("/", func(r chi.Router) {
				r.Put("/", api.updateProblemByName)
				r.Delete("/", api.deleteProblem)

				r.Post("/publish", api.publishProblem)
				r.Post("/unpublish", api.unpublishProblem)

				r.Route("/tests", func(r chi.Router) {
					r.Get("/", api.getProblemTests)
					r.Post("/", api.createProblemTest)

					r.With(api.problemTestCtx).Route("/{problemTestId}", func(r chi.Router) {
						r.Get("/", api.getProblemTestByID)
						r.Put("/", api.updateProblemTestById)
						r.Delete("/", api.deleteProblemTestById)
					})
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
