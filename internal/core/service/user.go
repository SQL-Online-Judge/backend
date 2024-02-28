package service

import (
	"fmt"

	"github.com/SQL-Online-Judge/backend/internal/core/repository"
	"github.com/SQL-Online-Judge/backend/internal/model"
)

var (
	ErrUserIsNil                 = fmt.Errorf("user is nil")
	ErrInvalidUsernameOrPassword = fmt.Errorf("invalid username or password")
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(ur repository.UserRepository) *UserService {
	return &UserService{
		repo: ur,
	}
}

func (us *UserService) CreateUser(u *model.User) error {
	err := us.repo.CreateUser(u)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (us *UserService) CheckPassword(u *model.User) error {
	if u == nil {
		return fmt.Errorf("%w", ErrUserIsNil)
	}

	if u.Username == "" || u.Password == "" {
		return fmt.Errorf("%w", ErrInvalidUsernameOrPassword)
	}

	user, err := us.repo.GetUserByUsername(u.Username)
	if err != nil {
		return fmt.Errorf("failed to get user by username: %w", err)
	}

	err = user.ComparePassword(u.Password)
	if err != nil {
		return fmt.Errorf("failed to compare password: %w", err)
	}

	return nil
}

func (us *UserService) GetUserIDByUsername(username string) int64 {
	user, err := us.repo.GetUserByUsername(username)
	if err != nil {
		return 0
	}
	return user.UserID
}
