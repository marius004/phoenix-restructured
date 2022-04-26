package repositories

import (
	"errors"
	"strings"

	"github.com/marius004/phoenix/entities"
	"github.com/marius004/phoenix/internal"
	"github.com/marius004/phoenix/models"
	"gorm.io/gorm"
)

type ProblemRepository struct {
	db *internal.Database
}

func (r *ProblemRepository) CreateProblem(problem *entities.Problem) error {
	result := r.db.Conn.Create(&problem)
	return result.Error
}

func (r *ProblemRepository) GetProblemByID(id uint) (*entities.Problem, error) {
	var problem *entities.Problem
	result := r.db.Conn.Where("id = ?", id).First(&problem)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return problem, result.Error
}

func (r *ProblemRepository) GetProblemByName(name string) (*entities.Problem, error) {
	var problem *entities.Problem
	result := r.db.Conn.Where("name = ?", name).First(&problem)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return problem, result.Error
}

func (r *ProblemRepository) GetProblemsByAuthorID(authorId uint) ([]*entities.Problem, error) {
	var problems []*entities.Problem
	result := r.db.Conn.Where("author_id = ?", authorId).Find(&problems)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return problems, result.Error
}

func (r *ProblemRepository) DeleteProblem(problem *entities.Problem) error {
	result := r.db.Conn.Unscoped().Delete(&entities.Problem{}, "id = ?", problem.ID)
	return result.Error
}

func (r *ProblemRepository) UpdateProblemByID(id uint, request *models.UpdateProblemRequest) error {
	problem, err := r.GetProblemByID(id)

	if err != nil {
		return err
	}

	if problem == nil {
		return internal.ErrProblemDoesNotExist
	}

	if request.Description != "" {
		problem.Description = request.Description
	}

	if request.MemoryLimit > 0 {
		problem.MemoryLimit = request.MemoryLimit
	}

	if request.StackLimit > 0 {
		problem.StackLimit = request.StackLimit
	}

	if request.TimeLimit > 0 {
		problem.TimeLimit = request.TimeLimit
	}

	if problem.Difficulty != "" {
		problem.Difficulty = request.Difficulty
	}

	result := r.db.Conn.Save(&problem)
	return result.Error
}

func (r *ProblemRepository) GetProblemsByFilter(filter *models.ProblemFilter) ([]*entities.Problem, error) {
	var problems []*entities.Problem
	query, args := makeProblemFilter(filter)

	var result *gorm.DB

	if filter.Limit > 0 && filter.Offset > 0 {
		result = r.db.Conn.Where(strings.Join(query, " AND "), args...).Offset(filter.Offset).Limit(filter.Limit).Find(&problems)
	} else if filter.Limit > 0 {
		result = r.db.Conn.Where(strings.Join(query, " AND "), args...).Limit(filter.Limit).Find(&problems)
	} else if filter.Offset > 0 {
		result = r.db.Conn.Where(strings.Join(query, " AND "), args...).Offset(filter.Offset).Find(&problems)
	} else {
		result = r.db.Conn.Where(strings.Join(query, " AND "), args...).Find(&problems)
	}

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return problems, result.Error
}

func makeProblemFilter(filter *models.ProblemFilter) (query []string, args []interface{}) {
	if filter.AuthorId > 0 {
		query = append(query, "author_id = ?")
		args = append(args, filter.AuthorId)
	}

	if filter.ProblemId > 0 {
		query = append(query, "id = ?")
		args = append(args, filter.ProblemId)
	}

	return
}

func NewProblemRepository(db *internal.Database) *ProblemRepository {
	return &ProblemRepository{
		db: db,
	}
}
