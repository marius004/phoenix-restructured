package repositories

import (
	"errors"

	"github.com/marius004/phoenix/entities"
	"github.com/marius004/phoenix/internal"
	"gorm.io/gorm"
)

type SubmissionRepository struct {
	db *internal.Database
}

func (r *SubmissionRepository) CreateSubmission(submission *entities.Submission) error {
	result := r.db.Conn.Create(&submission)
	return result.Error
}

func (r *SubmissionRepository) GetSubmissionByID(submissionId uint) (*entities.Submission, error) {
	var submission *entities.Submission
	result := r.db.Conn.Where("id = ?", submissionId).First(&submission)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return submission, result.Error
}

func (r *SubmissionRepository) GetSubmissionByUserID(userId uint) (*entities.Submission, error) {
	var submission *entities.Submission
	result := r.db.Conn.Where("user_id = ?", userId).First(&submission)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return submission, result.Error
}

func (r *SubmissionRepository) GetSubmissionByUsername(username string) (*entities.Submission, error) {
	var submission *entities.Submission
	result := r.db.Conn.Where("username = ?", username).First(&submission)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return submission, result.Error
}

func (r *SubmissionRepository) GetAllSubmissions() ([]*entities.Submission, error) {
	var submissions []*entities.Submission
	result := r.db.Conn.Find(&submissions)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return submissions, result.Error
}

func NewSubmissionRepository(db *internal.Database) *SubmissionRepository {
	return &SubmissionRepository{
		db: db,
	}
}
