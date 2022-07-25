package services

import (
	"context"
	"errors"
	"strings"

	"github.com/marius004/phoenix/entities"
	"github.com/marius004/phoenix/internal"
	"github.com/marius004/phoenix/models"
	"gorm.io/gorm"
)

type ProblemService struct {
	db *internal.Database
}

func (s *ProblemService) CreateProblem(context context.Context, problem *entities.Problem) error {
	if problem, _ := s.GetProblemByName(context, problem.Name); problem != nil {
		return internal.ErrProblemNameAlreadyExists
	}

	result := s.db.Conn.Create(&problem)
	return result.Error
}

func (s *ProblemService) GetProblemByID(context context.Context, id uint) (*entities.Problem, error) {
	var problem *entities.Problem
	result := s.db.Conn.Where("id = ?", id).First(&problem)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return problem, result.Error
}

func (s *ProblemService) GetProblemByName(context context.Context, name string) (*entities.Problem, error) {
	var problem *entities.Problem
	result := s.db.Conn.Where("name = ?", name).First(&problem)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return problem, result.Error
}

func (s *ProblemService) GetProblemsByAuthorID(context context.Context, authorId uint) ([]*entities.Problem, error) {
	var problems []*entities.Problem
	result := s.db.Conn.Where("author_id = ?", authorId).Find(&problems)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return problems, result.Error
}

func (s *ProblemService) GetProblemsByFilter(context context.Context, filter *models.ProblemFilter) ([]*entities.Problem, error) {
	var problems []*entities.Problem
	query, args := makeProblemFilter(filter)

	var result *gorm.DB

	if filter.Limit > 0 && filter.Offset > 0 {
		result = s.db.Conn.Where(strings.Join(query, " AND "), args...).Offset(filter.Offset).Limit(filter.Limit).Find(&problems)
	} else if filter.Limit > 0 {
		result = s.db.Conn.Where(strings.Join(query, " AND "), args...).Limit(filter.Limit).Find(&problems)
	} else if filter.Offset > 0 {
		result = s.db.Conn.Where(strings.Join(query, " AND "), args...).Offset(filter.Offset).Find(&problems)
	} else {
		result = s.db.Conn.Where(strings.Join(query, " AND "), args...).Find(&problems)
	}

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return problems, result.Error
}

func (s *ProblemService) UpdateProblemByID(context context.Context, id uint, user *entities.User, request *models.UpdateProblemRequest) error {
	problem, err := s.GetProblemByID(context, id)
	if problem == nil || err != nil {
		return internal.ErrProblemDoesNotExist
	}

	if !internal.IsUserProposer(user) ||
		(internal.IsUserProposer(user) && !internal.IsUserAdmin(user) && user.ID != problem.AuthorId) {
		return internal.ErrUnauthorized
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

	result := s.db.Conn.Save(&problem)
	return result.Error
}

func (s *ProblemService) DeleteProblem(context context.Context, problem *entities.Problem) error {
	result := s.db.Conn.Unscoped().Delete(&entities.Problem{}, "id = ?", problem.ID)
	return result.Error
}

func (s *ProblemService) UpdateProblemStatus(context context.Context, problem *entities.Problem, status entities.ProblemStatus) error {
	problem.Status = status

	result := s.db.Conn.Save(&problem)
	return result.Error
}

func NewProblemService(db *internal.Database) *ProblemService {
	return &ProblemService{
		db: db,
	}
}
