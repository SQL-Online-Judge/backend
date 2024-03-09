package service

import (
	"fmt"

	"github.com/SQL-Online-Judge/backend/internal/core/repository"
	"github.com/SQL-Online-Judge/backend/internal/model"
)

var (
	ErrUserIsNil                 = fmt.Errorf("user is nil")
	ErrInvalidUsernameOrPassword = fmt.Errorf("invalid username or password")
	ErrUserConflict              = fmt.Errorf("user conflict")
	ErrUserNotFound              = fmt.Errorf("user not found")
	ErrUserNotStudent            = fmt.Errorf("user is not student")
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(ur repository.UserRepository) *UserService {
	return &UserService{
		repo: ur,
	}
}

func (us *UserService) isUserIDExist(userID int64) bool {
	return us.repo.ExistByUserID(userID)
}

func (us *UserService) isUsernameExist(username string) bool {
	return us.repo.ExistByUsername(username)
}

func (us *UserService) isUserDeleted(userID int64) bool {
	return us.repo.IsDeletedByUserID(userID)
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
	if !us.isUsernameExist(username) {
		return 0, fmt.Errorf("%w", ErrUserNotFound)
	}

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
	if !us.isUsernameExist(username) {
		return 0
	}

	user, err := us.repo.FindByUsername(username)
	if err != nil {
		return 0
	}
	return user.UserID
}

func (us *UserService) GetRoleByUserID(userID int64) string {
	if !us.isUserIDExist(userID) {
		return ""
	}

	user, err := us.repo.FindByUserID(userID)
	if err != nil {
		return ""
	}
	return user.Role
}

func (us *UserService) DeleteByUserID(userID int64) error {
	if !us.isUserIDExist(userID) {
		return fmt.Errorf("%w", ErrUserNotFound)
	}

	err := us.repo.DeleteByUserID(userID)
	if err != nil {
		return fmt.Errorf("failed to delete user by userID: %w", err)
	}

	return nil
}

func (us *UserService) isStudentExist(userID int64) error {
	if !us.isUserIDExist(userID) {
		return fmt.Errorf("%w", ErrUserNotFound)
	}

	if us.isUserDeleted(userID) {
		return fmt.Errorf("%w", ErrUserNotFound)
	}

	role, err := us.repo.GetRoleByUserID(userID)
	if err != nil {
		return fmt.Errorf("failed to get role by userID: %w", err)
	}

	if role != "student" {
		return fmt.Errorf("%w", ErrUserNotStudent)
	}

	return nil
}

func (us *UserService) UpdateStudentUsername(userID int64, username string) error {
	if err := us.isStudentExist(userID); err != nil {
		return err
	}

	if us.isUsernameExist(username) {
		return fmt.Errorf("%w", ErrUserConflict)
	}

	if err := us.repo.UpdateUsernameByUserID(userID, username); err != nil {
		return fmt.Errorf("failed to update username: %w", err)
	}

	return nil
}

func (us *UserService) GetStudent(userID int64) (*model.User, error) {
	if err := us.isStudentExist(userID); err != nil {
		return nil, err
	}

	user, err := us.repo.FindByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student: %w", err)
	}
	return user, nil
}

func (us *UserService) GetStudents(contains string) ([]*model.User, error) {
	students, err := us.repo.GetStudents(contains)
	if err != nil {
		return nil, fmt.Errorf("failed to get students: %w", err)
	}
	return students, nil
}
