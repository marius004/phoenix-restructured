package internal

import (
	"github.com/marius004/phoenix/entities"
	"github.com/marius004/phoenix/models"
)

type Repositories struct {
	UserRepository UserRepository

	ProblemRepository     ProblemRepository
	ProblemTestRepository ProblemTestRepository

	SubmissionRepository     SubmissionRepository
	SubmissionTestRepository SubmissionTestRepository

	PostRepository PostRepository
}

type UserRepository interface {
	CreateUser(user *entities.User) error

	GetUserByID(id uint) (*entities.User, error)
	GetUserByEmail(email string) (*entities.User, error)
	GetUserByUsername(username string) (*entities.User, error)
	GetUsers(filter *models.UserFilter) ([]*entities.User, error)

	UpdateUser(user *entities.User, request *models.UpdateUserRequest) error

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

type SubmissionRepository interface {
	CreateSubmission(submission *entities.Submission) error

	GetBySubmissionFilter(filter models.SubmissionFilter) ([]*entities.Submission, error)
	GetAllSubmissions() ([]*entities.Submission, error)

	GetSubmissionByID(submissionId uint) (*entities.Submission, error)
	GetSubmissionByUserID(userId uint) (*entities.Submission, error)
	GetSubmissionByUsername(username string) (*entities.Submission, error)

	UpdateSubmission(submissionId uint, request *models.UpdateSubmissionRequest) error
}

type SubmissionTestRepository interface {
	GetSubmissionTestsBySubmissionID(submissionId uint) ([]*entities.SubmissionTest, error)
	GetSubmissionTestByID(submissionTestId uint) (*entities.SubmissionTest, error)
	GetSubmissionTestByProblemTestAndSubmissionID(problemTestId, submissionId uint) (*entities.SubmissionTest, error)

	CreateSubmissionTest(submissionTest *entities.SubmissionTest) error
	UpdateSubmissionTest(testId, submissionId uint, request *models.UpdateSubmissionTestRequest) error
}

type PostRepository interface {
	GetPosts() ([]*entities.Post, error)
	GetPostByTitle(title string) (*entities.Post, error)
	GetPostByID(id uint) (*entities.Post, error)

	CreatePost(post *entities.Post) error

	UpdatePostByID(id uint, request *models.UpdatePostRequest) error
	UpdatePostByTitle(title string, request *models.UpdatePostRequest) error

	DeletePost(post *entities.Post) error
}
