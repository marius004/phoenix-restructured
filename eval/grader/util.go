package grader

import (
	"context"

	"github.com/marius004/phoenix/entities"
	"github.com/marius004/phoenix/models"
)

var (
	ctx = context.Background()

	waitingSubmissions   = models.SubmissionFilter{Status: entities.Waiting}
	evaluatingSubmission = models.UpdateSubmissionRequest{Status: entities.Evaluating}
)
