package services

import (
	"context"
	"errors"

	"github.com/marius004/phoenix/entities"
	"github.com/marius004/phoenix/internal"
	"github.com/marius004/phoenix/models"
	"gorm.io/gorm"
)

type SubmissionTestService struct {
	db *internal.Database
}

func (s *SubmissionTestService) GetSubmissionTestsBySubmissionID(context context.Context, submissionId uint) ([]*entities.SubmissionTest, error) {
	var submissionTests []*entities.SubmissionTest
	result := s.db.Conn.Where("submission_id = ?", submissionId).Find(&submissionTests)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return submissionTests, result.Error
}

func (s *SubmissionTestService) GetSubmissionTestByID(context context.Context, submissionTestId uint) (*entities.SubmissionTest, error) {
	var submissionTest *entities.SubmissionTest
	result := s.db.Conn.Where("id = ?", submissionTestId).First(&submissionTest)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return submissionTest, result.Error
}

func (s *SubmissionTestService) GetSubmissionTestByTestAndSubmissionID(context context.Context, problemTestId, submissionId uint) (*entities.SubmissionTest, error) {
	var submissionTest *entities.SubmissionTest
	result := s.db.Conn.Where("submission_id = ? AND problem_test_id = ?", submissionId, problemTestId).First(&submissionTest)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return submissionTest, result.Error
}

func (s *SubmissionTestService) CreateSubmissionTest(context context.Context, submissionTest *entities.SubmissionTest) error {
	result := s.db.Conn.Create(&submissionTest)
	return result.Error
}

func (s *SubmissionTestService) UpdateSubmissionTest(context context.Context, testId, submissionId uint, request *models.UpdateSubmissionTestRequest) error {
	submissionTest, err := s.GetSubmissionTestByTestAndSubmissionID(context, testId, submissionId)

	if err != nil || submissionTest == nil {
		return internal.ErrSubmissionTestDoesNotExist
	}

	if request.ExecutionMessage != "" {
		submissionTest.ExecutionMessage = request.ExecutionMessage
	}

	if request.Score > 0 {
		submissionTest.Score = request.Score
	}

	if request.Time > 0 {
		submissionTest.Time = request.Time
	}

	if request.Memory > 0 {
		submissionTest.Memory = request.Memory
	}

	if request.ExitCode != 0 {
		submissionTest.ExitCode = request.ExitCode
	}

	result := s.db.Conn.Save(&submissionTest)
	return result.Error
}

func NewSubmissionTestService(db *internal.Database) *SubmissionTestService {
	return &SubmissionTestService{db}
}
