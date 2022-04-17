package repositories

import (
	"errors"

	"github.com/marius004/phoenix/entities"
	"github.com/marius004/phoenix/internal"
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

func (r *UserRepository) DeleteUser(user *entities.User) error {
	result := r.db.Conn.Delete(&user)
	return result.Error
}

func NewUserRepository(db *internal.Database) *UserRepository {
	return &UserRepository{
		db: db,
	}
}
