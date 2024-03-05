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
}
