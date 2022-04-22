package internal

import (
	"context"
	"encoding/json"
	"io"
	"io/fs"
	"os"

	"github.com/marius004/phoenix/entities"
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

	ExecuteCommand(ctx context.Context, command []string, config *RunConfig) (*RunStatus, error)
	Cleanup() error
}

// Task represents a task executed within a sandbox
type Task interface {
	Run(ctx context.Context, sandbox Sandbox) error
}

// RunConfig represents the configuration needed for a task to be run.
type RunConfig struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer

	MemoryLimit int
	StackLimit  int

	InputPath  string
	OutputPath string

	TimeLimit     float64
	WallTimeLimit float64

	MaxProcesses int
}

// RunStatus contains information about the process that was run within the sandbox
type RunStatus struct {
	Memory int

	ExitCode   int
	ExitSignal int
	Killed     bool

	Message string
	Status  string

	Time     float64
	WallTime float64
}

type CompileRequest struct {
	ID         uint
	SourceCode []byte
	Lang       string
}

type CompileResponse struct {
	Message string
	Success bool
}

type Limit struct {
	Time   float64
	Memory int
	Stack  int
}

type ExecuteRequest struct {
	ID uint

	SubmissionId int
	TestId       int

	Limit

	Lang      string
	ProblemId uint

	Input []byte
}

type ExecuteResponse struct {
	TimeUsed   float64
	MemoryUsed int

	ExitCode int
	Message  string
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

type EvaluatorServices struct {
	ProblemService     ProblemService
	ProblemTestService ProblemTestService

	SubmissionService     SubmissionService
	SubmissionTestService SubmissionTestService
}

func NewEvaluatorServices(services *Services) *EvaluatorServices {
	return &EvaluatorServices{
		ProblemService:     services.ProblemService,
		ProblemTestService: services.ProblemTestService,

		SubmissionService:     services.SubmissionService,
		SubmissionTestService: services.SubmissionTestService,
	}
}
