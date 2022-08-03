package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/marius004/phoenix-algo/api"
	"github.com/marius004/phoenix-algo/entities"
	"github.com/marius004/phoenix-algo/internal"
	"github.com/marius004/phoenix-algo/services"
	"github.com/marius004/phoenix-algo/services/eval/grader"
)

type Server struct {
	db         *internal.Database
	config     *internal.Config
	evalConfig *internal.EvalConfig

	services *internal.Services
}

var allEntities []interface{} = []interface{}{
	&entities.User{},

	&entities.Problem{},
	&entities.ProblemTest{},

	&entities.Submission{},
	&entities.SubmissionTest{},
}

func (s *Server) Serve() {
	var api = api.NewAPI(s.config, s.services)
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
	return &Server{
		db:         db,
		config:     config,
		evalConfig: evalConfig,

		services: createServices(db, config, evalConfig),
	}
}

func createServices(db *internal.Database, config *internal.Config, evalConfig *internal.EvalConfig) *internal.Services {
	var (
		submissionService     = services.NewSubmissionService(db)
		submissionTestService = services.NewSubmissionTestService(db)

		problemService     = services.NewProblemService(db)
		problemTestService = services.NewProblemTestService(db)

		userService = services.NewUserService(db, submissionService, problemService)

		graderServices = internal.NewGraderServices(problemService, problemTestService, submissionService, submissionTestService)
	)

	return &internal.Services{
		UserService: userService,

		ProblemService:     problemService,
		ProblemTestService: problemTestService,

		SubmissionService:     submissionService,
		SubmissionTestService: submissionTestService,

		Grader: grader.NewGrader(300*time.Millisecond, graderServices, evalConfig),
	}
}
