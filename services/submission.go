package services

import (
	"github.com/marius004/phoenix/entities"
	"github.com/marius004/phoenix/internal"
)

type SubmissionService struct {
	submissionRepository internal.SubmissionRepository
}

func (s *SubmissionService) CreateSubmission(submission *entities.Submission) error {
	return s.submissionRepository.CreateSubmission(submission)
}

func (s *SubmissionService) GetSubmissionByID(submissionId uint) (*entities.Submission, error) {
	return s.submissionRepository.GetSubmissionByID(submissionId)
}

func (s *SubmissionService) GetSubmissionByUserID(userId uint) (*entities.Submission, error) {
	return s.submissionRepository.GetSubmissionByUserID(userId)
}

func (s *SubmissionService) GetSubmissionByUsername(username string) (*entities.Submission, error) {
	return s.submissionRepository.GetSubmissionByUsername(username)
}

func (s *SubmissionService) GetAllSubmissions() ([]*entities.Submission, error) {
	return s.submissionRepository.GetAllSubmissions()
}

func NewSubmissionService(submissionRepository internal.SubmissionRepository) *SubmissionService {
	return &SubmissionService{submissionRepository}
}
