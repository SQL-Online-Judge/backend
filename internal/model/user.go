package model

import (
	"fmt"

	"github.com/SQL-Online-Judge/backend/internal/pkg/id"
	"github.com/SQL-Online-Judge/backend/internal/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidPassword = fmt.Errorf("invalid password")

type User struct {
	UserID           int64  `bson:"userID"`
	Role             string `bson:"role"`
	Username         string `bson:"username"`
	Password         string `bson:"password"`
	IsPasswordHashed bool   `bson:"-"`
}

func NewEmptyUser() *User {
	return &User{
		UserID: id.NewID(),
	}
}

func (u *User) IsValidPassword() bool {
	return len(u.Password) >= 8
}

func (u *User) HashPassword() error {
	if u.IsPasswordHashed {
		logger.Logger.Warn("password already hashed", zap.String("user", u.Username))
		return nil
	}

	if !u.IsValidPassword() {
		logger.Logger.Error("invalid password", zap.String("user", u.Username))
		return fmt.Errorf("invalid password: %w", ErrInvalidPassword)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Logger.Error("failed to hash password", zap.String("user", u.Username), zap.Error(err))
		return fmt.Errorf("failed to hash password: %w", err)
	}

	u.Password = string(hash)
	u.IsPasswordHashed = true
	logger.Logger.Info("successfully hashed password", zap.String("user", u.Username))
	return nil
}
