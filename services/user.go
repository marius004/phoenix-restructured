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

type UserService struct {
	db *internal.Database
}

func (s *UserService) CreateUser(context context.Context, user *entities.User) error {
	if usr, _ := s.GetUserByUsername(context, user.Username); usr != nil {
		return internal.ErrUsernameAlreadyExists
	}

	if usr, _ := s.GetUserByEmail(context, user.Email); usr != nil {
		return internal.ErrEmailAlreadyExists
	}

	passwordHash, err := internal.GeneratePasswordHash(user.Password)
	if err != nil {
		return internal.ErrCouldNotGeneratePasswordHash
	}

	user.Password = passwordHash

	result := s.db.Conn.Create(&user)
	return result.Error
}

func (s *UserService) GetUserByID(context context.Context, id uint) (*entities.User, error) {
	var user *entities.User
	result := s.db.Conn.Where("id = ?", id).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return user, result.Error
}

func (s *UserService) GetUserByEmail(context context.Context, email string) (*entities.User, error) {
	var user *entities.User
	result := s.db.Conn.Where("email = ?", email).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return user, result.Error
}

func (s *UserService) GetUserByUsername(context context.Context, username string) (*entities.User, error) {
	var user *entities.User
	result := s.db.Conn.Where("username = ?", username).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return user, result.Error
}

func (s *UserService) GetUsers(context context.Context, filter *models.UserFilter) ([]*entities.User, error) {
	var users []*entities.User
	query, args := makeUserFilter(filter)

	var result *gorm.DB

	if filter.Limit > 0 && filter.Offset > 0 {
		result = s.db.Conn.Where(strings.Join(query, " AND "), args...).Offset(filter.Offset).Limit(filter.Limit).Find(&users)
	} else if filter.Limit > 0 {
		result = s.db.Conn.Where(strings.Join(query, " AND "), args...).Limit(filter.Limit).Find(&users)
	} else if filter.Offset > 0 {
		result = s.db.Conn.Where(strings.Join(query, " AND "), args...).Offset(filter.Offset).Find(&users)
	} else {
		result = s.db.Conn.Where(strings.Join(query, " AND "), args...).Find(&users)
	}

	return users, result.Error
}

func (s *UserService) UpdateUser(context context.Context, user *entities.User, request *models.UpdateUserRequest) error {
	if request.Bio != "" {
		user.Bio = request.Bio
	}

	if request.GithubURL != "" {
		user.GithubURL = request.GithubURL
	}

	if request.LinkedInURL != "" {
		user.LinkedInURL = request.LinkedInURL
	}

	if request.WebsiteURL != "" {
		user.WebsiteURL = request.WebsiteURL
	}

	if request.UserIconURL != "" {
		user.UserIconURL = request.UserIconURL
	}

	result := s.db.Conn.Save(&user)
	return result.Error
}

func (s *UserService) DeleteUser(context context.Context, user *entities.User) error {
	result := s.db.Conn.Unscoped().Delete(&user)
	return result.Error
}

func (s *UserService) AssignProposerRole(context context.Context, username string, action bool) error {
	user, err := s.GetUserByUsername(context, username)
	if err != nil {
		return err
	}

	user.IsProposer = action

	result := s.db.Conn.Save(&user)
	return result.Error
}

func NewUserService(db *internal.Database) *UserService {
	return &UserService{db}
}
