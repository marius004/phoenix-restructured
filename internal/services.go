package internal

import (
	"context"

	"github.com/marius004/phoenix/entities"
	"github.com/marius004/phoenix/models"
)

type Services struct {
	UserService           UserService
	ProblemService        ProblemService
	ProblemTestService    ProblemTestService
	SubmissionService     SubmissionService
	SubmissionTestService SubmissionTestService

	Grader Grader
}

type UserService interface {
	CreateUser(context context.Context, user *entities.User) error

	GetUserByID(context context.Context, id uint) (*entities.User, error)
	GetUserByEmail(context context.Context, email string) (*entities.User, error)
	GetUserByUsername(context context.Context, username string) (*entities.User, error)
	GetUsers(context context.Context, filter *models.UserFilter) ([]*entities.User, error)

	UpdateUser(context context.Context, user *entities.User, request *models.UpdateUserRequest) error

	DeleteUser(context context.Context, user *entities.User) error
}

type ProblemService interface {
	CreateProblem(context context.Context, problem *entities.Problem) error

	GetProblemByID(context context.Context, id uint) (*entities.Problem, error)
	GetProblemByName(context context.Context, name string) (*entities.Problem, error)
	GetProblemsByAuthorID(context context.Context, authorId uint) ([]*entities.Problem, error)
	GetProblemsByFilter(context context.Context, filter *models.ProblemFilter) ([]*entities.Problem, error)

	UpdateProblemByID(context context.Context, id uint, user *entities.User, request *models.UpdateProblemRequest) error

	DeleteProblem(context context.Context, problem *entities.Problem) error
}

type ProblemTestService interface {
	CreateProblemTest(context context.Context, problemTest *entities.ProblemTest) error

	GetProblemTestByID(context context.Context, testId uint) (*entities.ProblemTest, error)
	GetProblemTestsByProblemID(context context.Context, problemId uint) ([]*entities.ProblemTest, error)

	UpdateProblemTestByID(context context.Context, testId uint, request *models.UpdateProblemTestRequest) error

	DeleteProblemTestByID(context context.Context, testId uint) error
	DeleteProblemTestByProblemID(context context.Context, problemId uint) error
}

type SubmissionService interface {
	CreateSubmission(context context.Context, submission *entities.Submission) error

	GetBySubmissionFilter(context context.Context, filter models.SubmissionFilter) ([]*entities.Submission, error)
	GetAllSubmissions(context context.Context) ([]*entities.Submission, error)

	GetSubmissionByID(context context.Context, submissionId uint) (*entities.Submission, error)
	GetSubmissionByUserID(context context.Context, userId uint) (*entities.Submission, error)

	UpdateSubmission(context context.Context, submissionId uint, request *models.UpdateSubmissionRequest) error
}

type SubmissionTestService interface {
	GetSubmissionTestsBySubmissionID(context context.Context, submissionId uint) ([]*entities.SubmissionTest, error)
	GetSubmissionTestByID(context context.Context, submissionTestId uint) (*entities.SubmissionTest, error)
	GetSubmissionTestByTestAndSubmissionID(context context.Context, testId, submissionId uint) (*entities.SubmissionTest, error)

	CreateSubmissionTest(context context.Context, submissionTest *entities.SubmissionTest) error
	UpdateSubmissionTest(context context.Context, testId, submissionId uint, request *models.UpdateSubmissionTestRequest) error
}

type Grader interface {
	Handle()
}
