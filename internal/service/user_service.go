package service

import (
	"context"
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
	Register(ctx context.Context, user *entity.User) error
	Login(ctx context.Context, email string, password string) (*entity.User, error)
	GetByID(ctx context.Context, id string) (*entity.User, error)
	GetByIDs(ctx context.Context, ids []string) ([]*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) Register(ctx context.Context, user *entity.User) error {
	existingUser, _ := s.userRepo.GetByEmail(ctx, user.Email)
	if existingUser != nil {
		return ErrEmailExists
	}

	user.ID = uuid.NewString()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	return s.userRepo.Create(ctx, user)
}

func (s *userService) Login(ctx context.Context, email string, password string) (*entity.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrInvalidPassword
	}

	return user, nil
}

func (s *userService) GetByID(ctx context.Context, id string) (*entity.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *userService) GetByIDs(ctx context.Context, ids []string) ([]*entity.User, error) {
	return s.userRepo.GetByIDs(ctx, ids)
}

func (s *userService) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}
