package services

import (
	"github.com/marius004/phoenix/entities"
	"github.com/marius004/phoenix/internal"
	"github.com/marius004/phoenix/models"
)

type UserService struct {
	userRepository internal.UserRepository
}

func (s *UserService) CreateUser(user *entities.User) error {
	if usr, _ := s.userRepository.GetUserByUsername(user.Username); usr != nil {
		return internal.ErrUsernameAlreadyExists
	}

	if usr, _ := s.userRepository.GetUserByEmail(user.Email); usr != nil {
		return internal.ErrEmailAlreadyExists
	}

	err := s.userRepository.CreateUser(user)
	return err
}

func (s *UserService) GetUserByID(id uint) (*entities.User, error) {
	return s.userRepository.GetUserByID(id)
}

func (s *UserService) GetUserByEmail(email string) (*entities.User, error) {
	return s.userRepository.GetUserByEmail(email)
}

func (s *UserService) GetUserByUsername(username string) (*entities.User, error) {
	return s.userRepository.GetUserByUsername(username)
}

func (s *UserService) UpdateUser(user *entities.User, request *models.UpdateUserRequest) error {
	return s.userRepository.UpdateUser(user, request)
}

func (s *UserService) DeleteUser(user *entities.User) error {
	return s.userRepository.DeleteUser(user)
}

func NewUserService(userRepository internal.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}
