package grader

import (
	"os"
	"sync"
	"time"

	"github.com/marius004/phoenix-/entities"
	"github.com/marius004/phoenix-/internal"
	"github.com/marius004/phoenix-/models"
	"github.com/marius004/phoenix-/services/eval/checker"
	"github.com/marius004/phoenix-/services/eval/sandbox"
	"github.com/marius004/phoenix-/services/eval/tasks"
)

type Grader struct {
	iterationInterval time.Duration

	services *internal.GraderServices
	manager  *sandbox.Manager

	evalConfig *internal.EvalConfig
}

func (g *Grader) Handle() {
	logger := internal.GetGlobalLoggerInstance()
	ticker := time.NewTicker(g.iterationInterval)

	for range ticker.C {
		submissions, err := g.services.SubmissionService.GetBySubmissionFilter(internal.DefaultCtx, waitingSubmissions)

		if err != nil {
			logger.Println(err)
			continue
		}

		if len(submissions) == 0 {
			continue
		}

		for _, submission := range submissions {
			if err := g.services.SubmissionService.UpdateSubmission(internal.DefaultCtx, submission.ID, &evaluatingSubmission); err != nil {
				logger.Println(err)
				continue
			}

			g.handleSubmission(submission)
		}
	}
}

func (g *Grader) handleSubmission(submission *entities.Submission) {
	logger := internal.GetGlobalLoggerInstance()

	if !g.compileSubmission(submission) {
		logger.Printf("Could not compile submission %d\n", submission.ID)
		return
	}

	if !g.executeSubmission(submission) {
		logger.Printf("Could not execute submission %d\n", submission.ID)
		return
	}

	checker := g.getAppropriateChecker()
	if err := checker.Check(submission); err != nil {
		logger.Printf("Could not check the submission %d\n", submission.ID)
		return
	}
}

func (g *Grader) compileSubmission(submission *entities.Submission) bool {
	logger := internal.GetGlobalLoggerInstance()
	logger.Printf("compiling submission %d\n", submission.ID)

	compile := &tasks.CompileTask{
		EvalConfig: g.evalConfig,
		Request: &models.CompileRequest{
			ID:         submission.ID,
			Lang:       string(submission.Language),
			SourceCode: submission.SourceCode,
		},
		Response: &models.CompileResponse{},
	}

	// try to compile
	if err := g.manager.RunTask(ctx, compile); err != nil {
		logger.Println(err)
		goto updateSubmission
	}

updateSubmission:

	updateSubmissionRequest := &models.UpdateSubmissionRequest{
		CompiledSuccesfully: &compile.Response.Success,
		Message:             compile.Response.Message,
	}

	if err := g.services.SubmissionService.UpdateSubmission(internal.DefaultCtx, submission.ID, updateSubmissionRequest); err != nil {
		logger.Println(err)
		return false
	}

	return compile.Response.Success
}

func (g *Grader) handleGraderError(submission *entities.Submission, message string) {
	updateSubmissionRequest := &models.UpdateSubmissionRequest{Message: message}

	if err := g.services.SubmissionService.UpdateSubmission(internal.DefaultCtx, submission.ID, updateSubmissionRequest); err != nil {
		internal.GetGlobalLoggerInstance().Println(err)
	}
}

func (g *Grader) createSubmissionTest(submissionTest *entities.SubmissionTest) error {
	return g.services.SubmissionTestService.CreateSubmissionTest(internal.DefaultCtx, submissionTest)
}

func (g *Grader) executeSubmission(submission *entities.Submission) bool {
	logger := internal.GetGlobalLoggerInstance()
	problem, err := g.services.ProblemService.GetProblemByID(internal.DefaultCtx, submission.ProblemId)

	if err != nil {
		logger.Println("could not fetch problem", err)
		g.handleGraderError(submission, "could not fetch problem")
		return false
	}

	problemTests, err := g.services.ProblemTestService.GetProblemTestsByProblemID(internal.DefaultCtx, problem.ID)
	if err != nil {
		logger.Println("could not fetch problem tests", err)
		g.handleGraderError(submission, "could not fetch problem tests")
		return false
	}

	var wg sync.WaitGroup
	for _, problemTest := range problemTests {
		wg.Add(1)

		go func(problemTest entities.ProblemTest) {
			defer wg.Done()

			execute := g.executeProblemTest(submission, problem, &problemTest)

			if err := g.manager.RunTask(ctx, execute); err != nil {
				logger.Println(err)
				submissionTest := &entities.SubmissionTest{
					SubmissionId:     submission.ID,
					ProblemTestId:    problemTest.ProblemId,
					ExecutionMessage: "could not execute test",
				}

				if err := g.createSubmissionTest(submissionTest); err != nil {
					logger.Println(err)
				}

				return
			}

			submissionTest := &entities.SubmissionTest{
				Time:   execute.Response.TimeUsed,
				Memory: execute.Response.MemoryUsed,

				ExecutionMessage: execute.Response.Message,
				ExitCode:         execute.Response.ExitCode,

				SubmissionId:  submission.ID,
				ProblemTestId: problemTest.ID,
			}

			if err := g.createSubmissionTest(submissionTest); err != nil {
				logger.Println(err)
				return
			}

		}(*problemTest)
	}

	wg.Wait()

	return true
}

func (g *Grader) executeProblemTest(submission *entities.Submission, problem *entities.Problem, problemTest *entities.ProblemTest) *tasks.ExecuteTask {
	return &tasks.ExecuteTask{
		EvalConfig: g.evalConfig,
		Request: &models.ExecuteRequest{
			ID: submission.ID,

			SubmissionId: int(submission.ID),
			TestId:       int(problemTest.ID),

			Limit: models.Limit{
				Time:   float64(problem.TimeLimit),
				Memory: problem.MemoryLimit,
				Stack:  problem.StackLimit,
			},

			Input:     problemTest.Input,
			ProblemId: problem.ID,
			Lang:      string(submission.Language),
		},

		Response: &models.ExecuteResponse{},
	}
}

func (g *Grader) getAppropriateChecker() internal.Checker {
	return checker.NewChecker(g.evalConfig, g.services)
}

func NewGrader(iterationInterval time.Duration, services *internal.GraderServices, evalConfig *internal.EvalConfig) *Grader {
	manager := sandbox.NewManager(evalConfig)

	os.Mkdir(evalConfig.CompilePath, 0777)
	os.Mkdir(evalConfig.OutputPath, 0777)

	return &Grader{
		iterationInterval: iterationInterval,

		services:   services,
		manager:    manager,
		evalConfig: evalConfig,
	}
}
