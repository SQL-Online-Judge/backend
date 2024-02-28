package repository

import "github.com/SQL-Online-Judge/backend/internal/model"

type Repository interface {
	UserRepository
}

type UserRepository interface {
	CreateUser(user *model.User) error
	GetUserByUsername(username string) (*model.User, error)
}
