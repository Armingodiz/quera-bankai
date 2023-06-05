package userService

import (
	"bankai/models"
	"bankai/repository/userRepository"
	"bankai/utils"
)

type UserService interface {
	CreateUser(user *models.User) error
	GetUser(username string) (*models.User, error)
	DeleteUser(username string) error
}

func NewUserService(userRepo userRepository.UserRepository) UserService {
	return &userService{
		userRepository: userRepo,
	}
}

type userService struct {
	userRepository userRepository.UserRepository
}

func (s *userService) CreateUser(user *models.User) error {
	hashed, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashed
	return s.userRepository.CreateUser(user)
}

func (s *userService) GetUser(username string) (*models.User, error) {
	return s.userRepository.GetUserByUsername(username)
}

func (s *userService) DeleteUser(username string) error {
	return nil
}
