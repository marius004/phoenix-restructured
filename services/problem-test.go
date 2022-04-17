package services

import (
	"github.com/marius004/phoenix/entities"
	"github.com/marius004/phoenix/internal"
	"github.com/marius004/phoenix/models"
)

type ProblemTestService struct {
	problemTestRepository internal.ProblemTestRepository
}

func (s *ProblemTestService) CreateProblemTest(problemTest *entities.ProblemTest) error {
	if err := s.problemTestRepository.CreateProblemTest(problemTest); err != nil {
		return err
	}

	return nil
}

func (s *ProblemTestService) GetProblemTestByID(testId uint) (*entities.ProblemTest, error) {
	problemTest, err := s.problemTestRepository.GetProblemTestByID(testId)

	if err != nil {
		return nil, err
	}

	return problemTest, nil
}

func (s *ProblemTestService) GetProblemTestsByProblemID(problemId uint) ([]*entities.ProblemTest, error) {
	problemTests, err := s.problemTestRepository.GetProblemTestsByProblemID(problemId)

	if err != nil {
		return nil, err
	}

	return problemTests, nil
}

func (s *ProblemTestService) UpdateProblemTestByID(testId uint, request *models.UpdateProblemTestRequest) error {
	if err := s.problemTestRepository.UpdateProblemTestByID(testId, request); err != nil {
		return err
	}

	return nil
}

func (s *ProblemTestService) DeleteProblemTestByID(testId uint) error {
	if err := s.problemTestRepository.DeleteProblemTestByID(testId); err != nil {
		return err
	}

	return nil
}

func (s *ProblemTestService) DeleteProblemTestByProblemID(problemId uint) error {
	if err := s.problemTestRepository.DeleteProblemTestsByProblemID(problemId); err != nil {
		return err
	}

	return nil
}

func NewProblemTestService(problemTestRepository internal.ProblemTestRepository) *ProblemTestService {
	return &ProblemTestService{
		problemTestRepository: problemTestRepository,
	}
}
