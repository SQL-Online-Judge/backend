package repository

import "github.com/SQL-Online-Judge/backend/internal/model"

type Repository interface {
	UserRepository
}

type UserRepository interface {
	Create(username, password, role string) (int64, error)
	FindByUserID(userID int64) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
	ExistByUserID(userID int64) bool
	ExistByUsername(username string) bool
	DeleteByUserID(userID int64) error
	GetRoleByUserID(userID int64) (string, error)
	UpdateUsernameByUserID(userID int64, username string) error
	IsDeletedByUserID(userID int64) bool
	GetStudents(contains string) ([]*model.User, error)
}
