package internal

import (
	"context"
	"encoding/json"
	"io/fs"
	"os"

	"github.com/marius004/phoenix-/entities"
	"github.com/marius004/phoenix-/models"
)

// Sandbox is an isolated process where code is executed
type Sandbox interface {
	GetID() int
	GetPath(path string) string

	CreateDirectory(path string, perm fs.FileMode) error
	DeleteDirectory(path string) error

	FileExists(path string) bool
	CreateFile(path string, perm fs.FileMode) error
	WriteToFile(path string, data []byte, perm fs.FileMode) error
	ReadFile(path string) ([]byte, error)
	DeleteFile(path string) error

	ExecuteCommand(ctx context.Context, command []string, config *models.RunConfig) (*models.RunStatus, error)
	Cleanup() error
}

// Task represents a task executed within a sandbox
type Task interface {
	Run(ctx context.Context, sandbox Sandbox) error
}

type Checker interface {
	Check(submission *entities.Submission) error
}

type EvalConfigLanguage struct {
	Extension  string `json:"extension"`
	IsCompiled bool   `json:"isCompiled"`

	Compile []string `json:"compile"`
	Execute []string `json:"execute"`

	SourceFile string `json:"sourceFile"`
	Executable string `json:"executable"`
}

type EvalConfig struct {
	IsolatePath  string `json:"isolatePath"`
	MaxSandboxes int    `json:"maxSandboxes"`
	CompilePath  string `json:"compilePath"`
	OutputPath   string `json:"outputPath"`

	Languages map[string]EvalConfigLanguage `json:"languages"`
}

func NewEvalConfig(evalConfigPath string) *EvalConfig {
	logger := GetGlobalLoggerInstance()
	evalConfig := &EvalConfig{}

	file, err := os.Open(evalConfigPath)
	if err != nil {
		logger.Fatalln("Could not open the evaluator json config file", err)
	}

	defer file.Close()
	if err := json.NewDecoder(file).Decode(&evalConfig); err != nil {
		logger.Fatalln("Could not decode the json config file", err)
	}

	return evalConfig
}

type GraderServices struct {
	ProblemService     ProblemService
	ProblemTestService ProblemTestService

	SubmissionService     SubmissionService
	SubmissionTestService SubmissionTestService
}

func NewGraderServices(problemService ProblemService,
	problemTestService ProblemTestService,
	submissionService SubmissionService,
	submissionTestService SubmissionTestService) *GraderServices {

	return &GraderServices{
		ProblemService:     problemService,
		ProblemTestService: problemTestService,

		SubmissionService:     submissionService,
		SubmissionTestService: submissionTestService,
	}
}
