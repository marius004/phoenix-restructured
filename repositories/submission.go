package repositories

import (
	"errors"
	"strings"

	"github.com/marius004/phoenix/entities"
	"github.com/marius004/phoenix/internal"
	"github.com/marius004/phoenix/models"
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
	result := r.db.Conn.Order("id desc").Find(&submissions)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return submissions, result.Error
}

func (r *SubmissionRepository) GetBySubmissionFilter(filter models.SubmissionFilter) ([]*entities.Submission, error) {
	var submissions []*entities.Submission
	query, args := makeFilter(filter)

	if len(query) == 0 {
		return r.GetAllSubmissions()
	}

	var result *gorm.DB

	if filter.Limit == -1 && filter.Offset == -1 {
		result = r.db.Conn.Where(strings.Join(query, " AND "), args...).Find(&submissions)
	} else if filter.Limit >= 0 && filter.Offset >= 0 {
		result = r.db.Conn.Where(strings.Join(query, " AND "), args...).Limit(filter.Limit).Offset(filter.Offset).Find(&submissions)
	} else if filter.Limit >= 0 {
		result = r.db.Conn.Where(strings.Join(query, " AND "), args...).Limit(filter.Limit).Find(&submissions)
	} else {
		result = r.db.Conn.Where(strings.Join(query, " AND "), args...).Offset(filter.Limit).Find(&submissions)
	}

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return submissions, result.Error
}

func (r *SubmissionRepository) UpdateSubmission(submissionId uint, request *models.UpdateSubmissionRequest) error {
	submission, err := r.GetSubmissionByID(submissionId)

	if err != nil {
		return internal.ErrSubmissionDoesNotExist
	}

	submission.Score = request.Score
	submission.Status = request.Status
	submission.Message = request.Message
	submission.CompiledSuccesfully = request.CompiledSuccesfully

	result := r.db.Conn.Save(&submission)
	return result.Error
}

func makeFilter(filter models.SubmissionFilter) ([]string, []interface{}) {
	var query []string
	var args []interface{}

	if filter.UserId > 0 {
		query = append(query, "user_id = ?")
		args = append(args, filter.UserId)
	}

	if filter.ProblemId > 0 {
		query = append(query, "problem_id = ?")
		args = append(args, filter.ProblemId)
	}

	if filter.Score > 0 {
		query = append(query, "score = ?")
		args = append(args, filter.Score)
	}

	if filter.Status != "" {
		query = append(query, "status = ?")
		args = append(args, filter.Status)
	}

	if filter.CompiledSuccesfully != nil {
		query = append(query, "compiled_succesfully = ?")
		args = append(args, filter.CompiledSuccesfully)
	}

	if filter.Limit > 0 {
		query = append(query, "limit = ?")
		args = append(args, filter.Limit)
	}

	if filter.Offset > 0 {
		query = append(query, "offset = ?")
		args = append(args, filter.Offset)
	}

	return query, args
}

func NewSubmissionRepository(db *internal.Database) *SubmissionRepository {
	return &SubmissionRepository{
		db: db,
	}
}
