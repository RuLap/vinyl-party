package service

import (
	"errors"
	"vinyl-party/internal/entity"
	"vinyl-party/internal/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrEmailExists     = errors.New("email exists")
	ErrInvalidPassword = errors.New("invalid password")
)

type UserService interface {
	Register(user *entity.User) error
	Login(email string, password string) (*entity.User, error)
	FindByID(id string) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) Register(user *entity.User) error {
	existingUser, _ := s.userRepo.FindByEmail(user.Email)
	if existingUser != nil {
		return ErrEmailExists
	}

	user.ID = uuid.NewString()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	return s.userRepo.Create(user)
}

func (s *userService) Login(email string, password string) (*entity.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrInvalidPassword
	}

	return user, nil
}

func (s *userService) FindByID(id string) (*entity.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (s *userService) FindByEmail(email string) (*entity.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}
