package models

type UpdateSubmissionTestRequest struct {
	ExecutionMessage string
	Score            int
	Time             float64
	Memory           int
	ExitCode         int
}
