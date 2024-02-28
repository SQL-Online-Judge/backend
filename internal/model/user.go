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
	UserID           int64  `bson:"userID" json:"userID"`
	Role             string `bson:"role" json:"role"`
	Username         string `bson:"username" json:"username"`
	Password         string `bson:"password" json:"password"`
	IsPasswordHashed bool   `bson:"-" json:"-"`
}

func NewEmptyUser() *User {
	return &User{
		UserID: id.NewID(),
	}
}

func (u *User) IsValidUsername() bool {
	return len(u.Username) >= 2 && len(u.Username) <= 32
}

func (u *User) IsValidPassword() bool {
	return len(u.Password) >= 8 && len(u.Password) <= 64
}

func (u *User) IsValidLogin() bool {
	return u.IsValidUsername() && u.IsValidPassword()
}

func (u *User) IsValidRole() bool {
	switch u.Role {
	case "admin", "teacher", "student":
		return true
	default:
		return false
	}
}

func (u *User) IsValidUser() bool {
	return u.IsValidUsername() && u.IsValidPassword() && u.IsValidRole()
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

func (u *User) ComparePassword(password string) error {
	if !u.IsPasswordHashed {
		logger.Logger.Error("password not hashed", zap.String("user", u.Username))
		return fmt.Errorf("password not hashed: %w", ErrInvalidPassword)
	}

	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return fmt.Errorf("failed to compare password: %w", err)
	}

	return nil
}

func (u *User) ClearPassword() {
	u.Password = ""
	u.IsPasswordHashed = false
}
