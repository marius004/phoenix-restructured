package checker

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/marius004/phoenix/entities"
	"github.com/marius004/phoenix/internal"
	"github.com/marius004/phoenix/models"
	"github.com/marius004/phoenix/services/eval"
)

type Checker struct {
	services   *internal.EvaluatorServices
	evalConfig *internal.EvalConfig
}

func (c *Checker) handleSubmissionError(submission *entities.Submission, message string) {
	updateSubmissionRequest := &models.UpdateSubmissionRequest{
		Score:  0,
		Status: entities.Evaluated,

		Message:             message,
		CompiledSuccesfully: &models.CompilationSuccess,
	}

	if err := c.services.SubmissionService.UpdateSubmission(submission.ID, updateSubmissionRequest); err != nil {
		internal.GetGlobalLoggerInstance().Println(err)
	}
}

func (c *Checker) handleSubmissionTestError(submission *entities.Submission, submissionTest *entities.SubmissionTest, message string) {
	updateSubmissionTestRequest := &models.UpdateSubmissionTestRequest{
		Score:  0,
		Time:   0,
		Memory: 0,

		ExecutionMessage: message,
		ExitCode:         0,
	}

	if err := c.services.SubmissionTestService.UpdateSubmissionTest(submissionTest.ProblemTestId, submission.ID, updateSubmissionTestRequest); err != nil {
		internal.GetGlobalLoggerInstance().Println(err)
	}
}

func (c *Checker) handleConstraintLimit(submission *entities.Submission, submissionTest *entities.SubmissionTest, constraint string) {
	updateSubmissionTestRequest := &models.UpdateSubmissionTestRequest{
		ExecutionMessage: constraint,
	}

	if err := c.services.SubmissionTestService.UpdateSubmissionTest(submissionTest.ProblemTestId, submission.ID, updateSubmissionTestRequest); err != nil {
		internal.GetGlobalLoggerInstance().Println(err)
	}
}

func (c *Checker) getSubmissionOutput(submission *entities.Submission, problemTest *entities.ProblemTest) ([]byte, error) {
	path := eval.GetOutputFileName(c.evalConfig, submission, problemTest)
	file, err := os.OpenFile(path, os.O_RDONLY, 0664)

	if err != nil {
		return nil, err
	}

	defer file.Close()
	return ioutil.ReadAll(file)
}

func (c *Checker) parseString(str string) string {
	str = strings.TrimSpace(str)
	str = strings.Trim(str, "\n")
	return str
}

func (c *Checker) isValidAnswer(received, expected []byte) bool {
	receivedStr := c.parseString(string(received))
	expectedStr := c.parseString(string(expected))

	return receivedStr == expectedStr
}

func (c *Checker) Check(submission *entities.Submission) error {
	logger := internal.GetGlobalLoggerInstance()
	problem, err := c.services.ProblemService.GetProblemByID(submission.ProblemId)

	if err != nil {
		logger.Println("could not fetch problem", err)
		c.handleSubmissionError(submission, "could not fetch problem")
		return err
	}

	problemTests, err := c.services.ProblemTestService.GetProblemTestsByProblemID(problem.ID)
	if err != nil {
		logger.Println("could not fetch problem tests", err)
		c.handleSubmissionError(submission, "could not fetch problem tests")
		return err
	}

	totalScore := 0
	for _, problemTest := range problemTests {
		submissionTest, err := c.services.SubmissionTestService.GetSubmissionTestByTestAndSubmissionID(problemTest.ID, submission.ID)

		if err != nil {
			c.handleSubmissionTestError(submission, submissionTest, "could not get submission test")
			continue
		}

		if submissionTest.Time > float64(problem.TimeLimit) {
			c.handleConstraintLimit(submission, submissionTest, "time limit exceeded")
			continue
		}

		if submissionTest.Memory > problem.MemoryLimit {
			c.handleConstraintLimit(submission, submissionTest, "memory limit exceeded")
			continue
		}

		submissionOutput, err := c.getSubmissionOutput(submission, problemTest)
		if err != nil {
			logger.Println(err)
			c.handleSubmissionTestError(submission, submissionTest, "could not get the submission output file")
			continue
		}

		var message = "Correct Answer"
		var testScore = problemTest.Score

		if !c.isValidAnswer(submissionOutput, problemTest.Output) {
			message = "Wrong Answer"
			testScore = 0
		}

		totalScore += testScore

		updateSubmissionTest := &models.UpdateSubmissionTestRequest{
			Score:            testScore,
			ExecutionMessage: message,
		}

		if err := c.services.SubmissionTestService.UpdateSubmissionTest(problemTest.ID, submission.ID, updateSubmissionTest); err != nil {
			logger.Println(err)
		}
	}

	logger.Println("TOTAL SCORE", totalScore)

	updateSubmissionRequest := &models.UpdateSubmissionRequest{
		Score:  totalScore,
		Status: entities.Evaluated,
	}

	if err := c.services.SubmissionService.UpdateSubmission(submission.ID, updateSubmissionRequest); err != nil {
		logger.Println(err)
		return err
	}

	return nil
}

func NewChecker(evalConfig *internal.EvalConfig, services *internal.EvaluatorServices) *Checker {
	return &Checker{
		services:   services,
		evalConfig: evalConfig,
	}
}
