package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"github.com/marius004/phoenix/api"
	"github.com/marius004/phoenix/entities"
	"github.com/marius004/phoenix/internal"
	"github.com/marius004/phoenix/repositories"
	"github.com/marius004/phoenix/services"
)

type Server struct {
	db         *internal.Database
	config     *internal.Config
	evalConfig *internal.EvalConfig

	repositories *internal.Repositories
	services     *internal.Services
}

var allEntities []interface{} = []interface{}{
	&entities.User{},

	&entities.Problem{},
	&entities.ProblemTest{},

	&entities.Submission{},
	&entities.SubmissionTest{},

	&entities.Post{},
}

func (s *Server) Serve() {
	var api = api.NewAPI(s.config, s.services, s.evalConfig)
	s.db.AutoMigrate(allEntities...)

	r := chi.NewRouter()

	corsConfig := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // https://stackoverflow.com/questions/24687313/what-exactly-does-the-access-control-allow-credentials-header-do
		MaxAge:           7200, // https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Max-Age
	})

	r.Use(corsConfig.Handler)
	r.Mount("/api", api.Routes())

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", s.config.ServerHost, s.config.ServerPort),
		Handler: r,
	}

	fmt.Printf("Phoenix running on %s:%s\n", s.config.ServerHost, s.config.ServerPort)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}

func NewServer(db *internal.Database, config *internal.Config, evalConfig *internal.EvalConfig) *Server {

	repositories := createRepositories(db)

	return &Server{
		db:         db,
		config:     config,
		evalConfig: evalConfig,

		repositories: repositories,
		services:     createServices(repositories, config),
	}
}

func createRepositories(db *internal.Database) *internal.Repositories {
	return &internal.Repositories{
		UserRepository: repositories.NewUserRepository(db),

		ProblemRepository:     repositories.NewProblemRepository(db),
		ProblemTestRepository: repositories.NewProblemTestRepository(db),

		SubmissionRepository:     repositories.NewSubmissionRepository(db),
		SubmissionTestRepository: repositories.NewSubmissionTestRepository(db),
	}
}

func createServices(repos *internal.Repositories, config *internal.Config) *internal.Services {
	return &internal.Services{
		UserService: services.NewUserService(repos.UserRepository),

		ProblemService:     services.NewProblemService(repos.ProblemRepository),
		ProblemTestService: services.NewProblemTestService(repos.ProblemTestRepository),

		SubmissionService:     services.NewSubmissionService(repos.SubmissionRepository),
		SubmissionTestService: services.NewSubmissionTestService(repos.SubmissionTestRepository),
	}
}
