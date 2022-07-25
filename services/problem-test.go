package services

import (
	"context"
	"errors"

	"github.com/marius004/phoenix/entities"
	"github.com/marius004/phoenix/internal"
	"github.com/marius004/phoenix/models"
	"gorm.io/gorm"
)

type ProblemTestService struct {
	db *internal.Database
}

func (s *ProblemTestService) CreateProblemTest(context context.Context, problemTest *entities.ProblemTest) error {
	result := s.db.Conn.Create(&problemTest)
	return result.Error
}

func (s *ProblemTestService) GetProblemTestByID(context context.Context, testId uint) (*entities.ProblemTest, error) {
	var problemTest *entities.ProblemTest
	result := s.db.Conn.Where("id = ?", testId).First(&problemTest)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return problemTest, result.Error
}

func (s *ProblemTestService) GetProblemTestsByProblemID(context context.Context, problemId uint) ([]*entities.ProblemTest, error) {
	var problemTests []*entities.ProblemTest
	result := s.db.Conn.Where("problem_id = ?", problemId).Find(&problemTests)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return problemTests, result.Error
}

func (s *ProblemTestService) UpdateProblemTestByID(context context.Context, testId uint, request *models.UpdateProblemTestRequest) error {
	problemTest, err := s.GetProblemTestByID(context, testId)

	if err != nil {
		return err
	} else if problemTest == nil {
		return internal.ErrProblemTestDoesNotExist
	}

	problemTest.Score = request.Score
	problemTest.Input = request.Input
	problemTest.Output = request.Output

	result := s.db.Conn.Save(&problemTest)
	return result.Error
}

func (s *ProblemTestService) DeleteProblemTestByID(context context.Context, testId uint) error {
	result := s.db.Conn.Unscoped().Where("id = ?", testId).Delete(&entities.ProblemTest{})
	return result.Error
}

func (s *ProblemTestService) DeleteProblemTestByProblemID(context context.Context, problemId uint) error {
	result := s.db.Conn.Unscoped().Where("problem_id = ?", problemId).Delete(&entities.ProblemTest{})
	return result.Error
}

func NewProblemTestService(db *internal.Database) *ProblemTestService {
	return &ProblemTestService{
		db: db,
	}
}
