package services

import (
	"github.com/marius004/phoenix/entities"
	"github.com/marius004/phoenix/internal"
	"github.com/marius004/phoenix/models"
)

type SubmissionTestService struct {
	submissionTestRepository internal.SubmissionTestRepository
}

func (s *SubmissionTestService) GetSubmissionTestsBySubmissionID(submissionId uint) ([]*entities.SubmissionTest, error) {
	return s.submissionTestRepository.GetSubmissionTestsBySubmissionID(submissionId)
}

func (s *SubmissionTestService) GetSubmissionTestByID(submissionTestId uint) (*entities.SubmissionTest, error) {
	return s.submissionTestRepository.GetSubmissionTestByID(submissionTestId)
}

func (s *SubmissionTestService) GetSubmissionTestByTestAndSubmissionID(testId, submissionId uint) (*entities.SubmissionTest, error) {
	return s.submissionTestRepository.GetSubmissionTestByTestAndSubmissionID(testId, submissionId)
}

func (s *SubmissionTestService) CreateSubmissionTest(submissionTest *entities.SubmissionTest) error {
	return s.submissionTestRepository.CreateSubmissionTest(submissionTest)
}

func (s *SubmissionTestService) UpdateSubmissionTest(testId, submissionId uint, request *models.UpdateSubmissionTestRequest) error {
	return s.submissionTestRepository.UpdateSubmissionTest(testId, submissionId, request)
}

func NewSubmissionTestService(submissionTestRepository internal.SubmissionTestRepository) *SubmissionTestService {
	return &SubmissionTestService{submissionTestRepository}
}
