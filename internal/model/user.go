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
	UserID   int64  `bson:"userID" json:"userID"`
	Username string `bson:"username" json:"username"`
	Password string `bson:"password" json:"password"`
	Role     string `bson:"role" json:"role"`
	Deleted  bool   `bson:"deleted"`
}

func NewEmptyUser() *User {
	return &User{
		UserID: id.NewID(),
	}
}

func NewUser(username, password, role string) *User {
	return &User{
		UserID:   id.NewID(),
		Username: username,
		Password: password,
		Role:     role,
		Deleted:  false,
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
	return u.IsValidUsername() && u.IsValidPassword() && u.IsValidRole() && !u.Deleted
}

func (u *User) GetHashedPassword() (string, error) {
	if !u.IsValidPassword() {
		logger.Logger.Error("invalid password", zap.String("user", u.Username))
		return "", fmt.Errorf("invalid password: %w", ErrInvalidPassword)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Logger.Error("failed to hash password", zap.String("user", u.Username), zap.Error(err))
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	hashedPassword := string(hash)
	return hashedPassword, nil
}

func (u *User) ComparePassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return fmt.Errorf("failed to compare password: %w", err)
	}

	return nil
}

func (u *User) ClearPassword() {
	u.Password = ""
}
