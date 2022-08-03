package services

import (
	"context"
	"errors"
	"strings"

	"github.com/marius004/phoenix-algo/entities"
	"github.com/marius004/phoenix-algo/internal"
	"github.com/marius004/phoenix-algo/models"
	"gorm.io/gorm"
)

type SubmissionService struct {
	db *internal.Database
}

func (s *SubmissionService) CreateSubmission(context context.Context, submission *entities.Submission) error {
	result := s.db.Conn.Create(&submission)
	return result.Error
}

func (s *SubmissionService) GetSubmissionByID(context context.Context, submissionId uint) (*entities.Submission, error) {
	var submission *entities.Submission
	result := s.db.Conn.Where("id = ?", submissionId).First(&submission)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return submission, result.Error
}

func (s *SubmissionService) GetSubmissionsByUserID(context context.Context, userId uint) ([]*entities.Submission, error) {
	var submissions []*entities.Submission
	result := s.db.Conn.Where("user_id = ?", userId).Find(&submissions)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return submissions, result.Error
}

func (s *SubmissionService) GetAllSubmissions(context context.Context) ([]*entities.Submission, error) {
	var submissions []*entities.Submission
	result := s.db.Conn.Order("id desc").Find(&submissions)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return submissions, result.Error
}

func (s *SubmissionService) GetBySubmissionFilter(context context.Context, filter models.SubmissionFilter) ([]*entities.Submission, error) {
	var submissions []*entities.Submission
	query, args := makeSubmissionFilter(filter)

	var result *gorm.DB

	if filter.Limit <= 0 && filter.Offset <= 0 {
		result = s.db.Conn.Order("id desc").Where(strings.Join(query, " AND "), args...).Find(&submissions)
	} else if filter.Limit >= 0 && filter.Offset >= 0 {
		result = s.db.Conn.Order("id desc").Where(strings.Join(query, " AND "), args...).Offset(filter.Offset).Limit(filter.Limit).Find(&submissions)
	} else if filter.Limit >= 0 {
		result = s.db.Conn.Order("id desc").Where(strings.Join(query, " AND "), args...).Limit(filter.Limit).Find(&submissions)
	} else {
		result = s.db.Conn.Order("id desc").Where(strings.Join(query, " AND "), args...).Offset(filter.Offset).Find(&submissions)
	}

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return submissions, result.Error
}

func (s *SubmissionService) UpdateSubmission(context context.Context, submissionId uint, request *models.UpdateSubmissionRequest) error {
	submission, err := s.GetSubmissionByID(context, submissionId)

	if err != nil {
		return internal.ErrSubmissionDoesNotExist
	}

	submission.Score = request.Score
	submission.Status = request.Status
	submission.Message = request.Message
	submission.CompiledSuccesfully = request.CompiledSuccesfully

	result := s.db.Conn.Save(&submission)
	return result.Error
}

func NewSubmissionService(db *internal.Database) *SubmissionService {
	return &SubmissionService{db}
}
