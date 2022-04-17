package services

import (
	"github.com/marius004/phoenix/entities"
	"github.com/marius004/phoenix/internal"
	"github.com/marius004/phoenix/models"
)

type ProblemService struct {
	problemRepository internal.ProblemRepository
}

func (s *ProblemService) CreateProblem(problem *entities.Problem) error {
	if problem, _ := s.problemRepository.GetProblemByName(problem.Name); problem != nil {
		return internal.ErrProblemNameAlreadyExists
	}

	err := s.problemRepository.CreateProblem(problem)
	return err
}

func (s *ProblemService) GetProblemByID(id uint) (*entities.Problem, error) {
	return s.problemRepository.GetProblemByID(id)
}

func (s *ProblemService) GetProblemByName(name string) (*entities.Problem, error) {
	return s.problemRepository.GetProblemByName(name)
}

func (s *ProblemService) GetProblemsByAuthorID(authorId uint) ([]*entities.Problem, error) {
	return s.problemRepository.GetProblemsByAuthorID(authorId)
}

func (s *ProblemService) GetProblemsByFilter(filter *models.ProblemFilter) ([]*entities.Problem, error) {
	return s.problemRepository.GetProblemsByFilter(filter)
}

func (s *ProblemService) UpdateProblemByID(id uint, user *entities.User, request *models.UpdateProblemRequest) error {
	problem, err := s.problemRepository.GetProblemByID(id)
	if problem == nil || err != nil {
		return internal.ErrProblemDoesNotExist
	}

	if !internal.IsUserProposer(user) ||
		(internal.IsUserProposer(user) && !internal.IsUserAdmin(user) && user.ID != problem.AuthorId) {
		return internal.ErrUnauthorized
	}

	return s.problemRepository.UpdateProblemByID(id, request)
}

func (s *ProblemService) DeleteProblem(problem *entities.Problem) error {
	return s.problemRepository.DeleteProblem(problem)
}

func NewProblemService(problemRepository internal.ProblemRepository) *ProblemService {
	return &ProblemService{
		problemRepository: problemRepository,
	}
}
