package repositories

import (
	"errors"

	"github.com/marius004/phoenix/entities"
	"github.com/marius004/phoenix/internal"
	"github.com/marius004/phoenix/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *internal.Database
}

func (r *UserRepository) CreateUser(user *entities.User) error {
	passwordHash, err := internal.GeneratePasswordHash(user.Password)
	if err != nil {
		return internal.ErrCouldNotGeneratePasswordHash
	}

	user.Password = passwordHash

	result := r.db.Conn.Create(&user)
	return result.Error
}

func (r *UserRepository) GetUserByID(id uint) (*entities.User, error) {
	var user *entities.User
	result := r.db.Conn.Where("id = ?", id).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return user, result.Error
}

func (r *UserRepository) GetUserByEmail(email string) (*entities.User, error) {
	var user *entities.User
	result := r.db.Conn.Where("email = ?", email).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return user, result.Error
}

func (r *UserRepository) GetUserByUsername(username string) (*entities.User, error) {
	var user *entities.User
	result := r.db.Conn.Where("username = ?", username).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return user, result.Error
}

func (r *UserRepository) UpdateUser(user *entities.User, request *models.UpdateUserRequest) error {
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

	result := r.db.Conn.Save(&user)
	return result.Error
}

func (r *UserRepository) DeleteUser(user *entities.User) error {
	result := r.db.Conn.Unscoped().Delete(&user)
	return result.Error
}

func NewUserRepository(db *internal.Database) *UserRepository {
	return &UserRepository{
		db: db,
	}
}
