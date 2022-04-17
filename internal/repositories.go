package internal

import (
	"github.com/marius004/phoenix/entities"
	"github.com/marius004/phoenix/models"
)

type Repositories struct {
	UserRepository UserRepository

	ProblemRepository     ProblemRepository
	ProblemTestRepository ProblemTestRepository
	SubmissionRepository  SubmissionRepository
}

type UserRepository interface {
	CreateUser(user *entities.User) error

	GetUserByID(id uint) (*entities.User, error)
	GetUserByEmail(email string) (*entities.User, error)
	GetUserByUsername(username string) (*entities.User, error)

	DeleteUser(user *entities.User) error
}

type ProblemRepository interface {
	CreateProblem(problem *entities.Problem) error

	GetProblemByID(id uint) (*entities.Problem, error)
	GetProblemByName(name string) (*entities.Problem, error)

	GetProblemsByAuthorID(authorId uint) ([]*entities.Problem, error)
	GetProblemsByFilter(filter *models.ProblemFilter) ([]*entities.Problem, error)

	UpdateProblemByID(id uint, request *models.UpdateProblemRequest) error

	DeleteProblem(problem *entities.Problem) error
}

type ProblemTestRepository interface {
	CreateProblemTest(problemTest *entities.ProblemTest) error

	GetProblemTestByID(id uint) (*entities.ProblemTest, error)
	GetProblemTestsByProblemID(id uint) ([]*entities.ProblemTest, error)

	UpdateProblemTestByID(testId uint, request *models.UpdateProblemTestRequest) error

	DeleteProblemTestByID(testId uint) error
	DeleteProblemTestsByProblemID(problemId uint) error
}

// TODO GetByFilter() :) and update submission
type SubmissionRepository interface {
	CreateSubmission(submission *entities.Submission) error

	GetSubmissionByID(submissionId uint) (*entities.Submission, error)
	GetSubmissionByUserID(userId uint) (*entities.Submission, error)
	GetSubmissionByUsername(username string) (*entities.Submission, error)
	GetAllSubmissions() ([]*entities.Submission, error)
}
