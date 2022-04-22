package repositories

import (
	"errors"

	"github.com/marius004/phoenix/entities"
	"github.com/marius004/phoenix/internal"
	"github.com/marius004/phoenix/models"
	"gorm.io/gorm"
)

type SubmissionTestRepository struct {
	db *internal.Database
}

func (r *SubmissionTestRepository) GetSubmissionTestsBySubmissionID(submissionId uint) ([]*entities.SubmissionTest, error) {
	var submissionTests []*entities.SubmissionTest
	result := r.db.Conn.Where("submission_id = ?", submissionId).Find(&submissionTests)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return submissionTests, result.Error
}

func (r *SubmissionTestRepository) GetSubmissionTestByID(submissionTestId uint) (*entities.SubmissionTest, error) {
	var submissionTest *entities.SubmissionTest
	result := r.db.Conn.Where("id = ?", submissionTestId).First(&submissionTest)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return submissionTest, result.Error
}

func (r *SubmissionTestRepository) GetSubmissionTestByTestAndSubmissionID(testId, submissionId uint) (*entities.SubmissionTest, error) {
	var submissionTest *entities.SubmissionTest
	result := r.db.Conn.Where("submission_id = ? AND problem_test_id = ?", submissionId, testId).First(&submissionTest)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return submissionTest, result.Error
}

func (r *SubmissionTestRepository) CreateSubmissionTest(submissionTest *entities.SubmissionTest) error {
	result := r.db.Conn.Create(&submissionTest)
	return result.Error
}

func (r *SubmissionTestRepository) UpdateSubmissionTest(testId, submissionId uint, request *models.UpdateSubmissionTestRequest) error {
	submissionTest, err := r.GetSubmissionTestByTestAndSubmissionID(testId, submissionId)

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

	result := r.db.Conn.Save(&submissionTest)
	return result.Error
}

func NewSubmissionTestRepository(db *internal.Database) *SubmissionTestRepository {
	return &SubmissionTestRepository{
		db: db,
	}
}
