package internal

import (
	"github.com/marius004/phoenix/entities"
	"github.com/marius004/phoenix/models"
)

type Services struct {
	UserService UserService

	ProblemService        ProblemService
	ProblemTestService    ProblemTestService
	SubmissionService     SubmissionService
	SubmissionTestService SubmissionTestService
	PostService           PostService
}

type UserService interface {
	CreateUser(user *entities.User) error

	GetUserByID(id uint) (*entities.User, error)
	GetUserByEmail(email string) (*entities.User, error)
	GetUserByUsername(username string) (*entities.User, error)
	GetUsers(filter *models.UserFilter) ([]*entities.User, error)

	UpdateUser(user *entities.User, request *models.UpdateUserRequest) error

	DeleteUser(user *entities.User) error
}

type ProblemService interface {
	CreateProblem(problem *entities.Problem) error

	GetProblemByID(id uint) (*entities.Problem, error)
	GetProblemByName(name string) (*entities.Problem, error)
	GetProblemsByAuthorID(authorId uint) ([]*entities.Problem, error)
	GetProblemsByFilter(filter *models.ProblemFilter) ([]*entities.Problem, error)

	UpdateProblemByID(id uint, user *entities.User, request *models.UpdateProblemRequest) error

	DeleteProblem(problem *entities.Problem) error
}

type ProblemTestService interface {
	CreateProblemTest(problemTest *entities.ProblemTest) error

	GetProblemTestByID(testId uint) (*entities.ProblemTest, error)
	GetProblemTestsByProblemID(problemId uint) ([]*entities.ProblemTest, error)

	UpdateProblemTestByID(testId uint, request *models.UpdateProblemTestRequest) error

	DeleteProblemTestByID(testId uint) error
	DeleteProblemTestByProblemID(problemId uint) error
}

type SubmissionService interface {
	CreateSubmission(submission *entities.Submission) error

	GetBySubmissionFilter(filter models.SubmissionFilter) ([]*entities.Submission, error)
	GetAllSubmissions() ([]*entities.Submission, error)

	GetSubmissionByID(submissionId uint) (*entities.Submission, error)
	GetSubmissionByUserID(userId uint) (*entities.Submission, error)
	GetSubmissionByUsername(username string) (*entities.Submission, error)

	UpdateSubmission(submissionId uint, request *models.UpdateSubmissionRequest) error
}

type SubmissionTestService interface {
	GetSubmissionTestsBySubmissionID(submissionId uint) ([]*entities.SubmissionTest, error)
	GetSubmissionTestByID(submissionTestId uint) (*entities.SubmissionTest, error)
	GetSubmissionTestByTestAndSubmissionID(testId, submissionId uint) (*entities.SubmissionTest, error)

	CreateSubmissionTest(submissionTest *entities.SubmissionTest) error
	UpdateSubmissionTest(testId, submissionId uint, request *models.UpdateSubmissionTestRequest) error
}

type PostService interface {
	GetPosts() ([]*entities.Post, error)
	GetPostByTitle(title string) (*entities.Post, error)
	GetPostByID(id uint) (*entities.Post, error)

	CreatePost(post *entities.Post) error

	UpdatePostByID(id uint, request *models.UpdatePostRequest) error
	UpdatePostByTitle(title string, request *models.UpdatePostRequest) error

	DeletePost(post *entities.Post) error
}
