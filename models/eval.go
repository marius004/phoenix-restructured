package models

import "io"

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
