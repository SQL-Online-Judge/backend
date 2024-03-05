package service

import (
	"fmt"

	"github.com/SQL-Online-Judge/backend/internal/core/repository"
)

var (
	ErrUserIsNil                 = fmt.Errorf("user is nil")
	ErrUserConflict              = fmt.Errorf("user conflict")
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

func (us *UserService) isUsernameExist(username string) bool {
	return us.repo.ExistByUsername(username)
}

func (us *UserService) CreateUser(username, password, role string) (int64, error) {
	if us.isUsernameExist(username) {
		return 0, fmt.Errorf("%w", ErrUserConflict)
	}

	userID, err := us.repo.Create(username, password, role)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}
	return userID, nil
}

func (us *UserService) Login(username, password string) (int64, error) {
	user, err := us.repo.FindByUsername(username)
	if err != nil {
		return 0, fmt.Errorf("failed to get user by username: %w", err)
	}

	err = user.ComparePassword(password)
	if err != nil {
		return 0, fmt.Errorf("failed to compare password: %w", err)
	}

	return user.UserID, nil
}

func (us *UserService) GetUserIDByUsername(username string) int64 {
	user, err := us.repo.FindByUsername(username)
	if err != nil {
		return 0
	}
	return user.UserID
}

func (us *UserService) GetRoleByUserID(userID int64) string {
	user, err := us.repo.FindByUserID(userID)
	if err != nil {
		return ""
	}
	return user.Role
}
