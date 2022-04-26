package repositories

import (
	"errors"
	"fmt"
	"strings"

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

func (r *UserRepository) GetUsers(filter *models.UserFilter) ([]*entities.User, error) {
	var users []*entities.User
	query, args := makeUserFilter(filter)

	var result *gorm.DB

	if filter.Limit > 0 && filter.Offset > 0 {
		result = r.db.Conn.Where(strings.Join(query, " AND "), args...).Offset(filter.Offset).Limit(filter.Limit).Find(&users)
	} else if filter.Limit > 0 {
		result = r.db.Conn.Where(strings.Join(query, " AND "), args...).Limit(filter.Limit).Find(&users)
	} else if filter.Offset > 0 {
		result = r.db.Conn.Where(strings.Join(query, " AND "), args...).Offset(filter.Offset).Find(&users)
	} else {
		result = r.db.Conn.Where(strings.Join(query, " AND "), args...).Find(&users)
	}

	return users, result.Error
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

func makeUserFilter(filter *models.UserFilter) (query []string, args []interface{}) {
	if filter.Email != "" {
		query = append(query, "email = ?")
		args = append(args, filter.Email)
	}

	if filter.GithubURL != "" {
		fmt.Println(filter.GithubURL)
		query = append(query, "github_url = ?")
		args = append(args, filter.GithubURL)
	}

	if filter.LinkedInURL != "" {
		query = append(query, "linked_in_url = ?")
		args = append(args, filter.LinkedInURL)
	}

	if filter.UserIconURL != "" {
		query = append(query, "user_icon_url = ?")
		args = append(args, filter.UserIconURL)
	}

	if filter.Username != "" {
		query = append(query, "username = ?")
		args = append(args, filter.Username)
	}

	if filter.WebsiteURL != "" {
		query = append(query, "website_url = ?")
		args = append(args, filter.WebsiteURL)
	}

	if filter.IsAdmin != nil {
		query = append(query, "is_admin = ?")
		args = append(args, filter.IsAdmin)
	}

	if filter.IsProposer != nil {
		query = append(query, "is_proposer = ?")
		args = append(args, filter.IsProposer)
	}

	if filter.UserId > 0 {
		query = append(query, "id = ?")
		args = append(args, filter.UserId)
	}

	return
}

func NewUserRepository(db *internal.Database) *UserRepository {
	return &UserRepository{
		db: db,
	}
}
