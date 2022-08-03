package grader

import (
	"context"

	"github.com/marius004/phoenix-algo/entities"
	"github.com/marius004/phoenix-algo/models"
)

var (
	ctx = context.Background()

	waitingSubmissions = models.SubmissionFilter{
		UserId:              -1,
		ProblemId:           -1,
		Score:               -1,
		CompiledSuccesfully: nil,
		Status:              entities.Waiting,
	}
	evaluatingSubmission = models.UpdateSubmissionRequest{
		Status: entities.Evaluating,
	}
)
