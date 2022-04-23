package repositories

import (
	"errors"

	"github.com/marius004/phoenix/entities"
	"github.com/marius004/phoenix/internal"
	"github.com/marius004/phoenix/models"
	"gorm.io/gorm"
)

type ProblemTestRepository struct {
	db *internal.Database
}

func (r *ProblemTestRepository) CreateProblemTest(problemTest *entities.ProblemTest) error {
	result := r.db.Conn.Create(&problemTest)
	return result.Error
}

func (r *ProblemTestRepository) GetProblemTestByID(id uint) (*entities.ProblemTest, error) {
	var problemTest *entities.ProblemTest
	result := r.db.Conn.Where("id = ?", id).First(&problemTest)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return problemTest, result.Error
}

func (r *ProblemTestRepository) GetProblemTestsByProblemID(id uint) ([]*entities.ProblemTest, error) {
	var problemTests []*entities.ProblemTest
	result := r.db.Conn.Where("problem_id = ?", id).Find(&problemTests)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return problemTests, result.Error
}

func (r *ProblemTestRepository) UpdateProblemTestByID(testId uint, request *models.UpdateProblemTestRequest) error {
	problemTest, err := r.GetProblemTestByID(testId)

	if err != nil {
		return err
	} else if problemTest == nil {
		return internal.ErrProblemTestDoesNotExist
	}

	problemTest.Score = request.Score
	problemTest.Input = request.Input
	problemTest.Output = request.Output

	r.db.Conn.Save(&problemTest)
	return nil
}

func (r *ProblemTestRepository) DeleteProblemTestByID(testId uint) error {
	result := r.db.Conn.Unscoped().Where("id = ?", testId).Delete(&entities.ProblemTest{})
	return result.Error
}

func (r *ProblemTestRepository) DeleteProblemTestsByProblemID(problemId uint) error {
	result := r.db.Conn.Unscoped().Where("problem_id = ?", problemId).Delete(&entities.ProblemTest{})
	return result.Error
}

func NewProblemTestRepository(db *internal.Database) *ProblemTestRepository {
	return &ProblemTestRepository{
		db: db,
	}
}
