package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/SQL-Online-Judge/backend/internal/model"
	"github.com/SQL-Online-Judge/backend/internal/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

var ErrUserIsNil = fmt.Errorf("user is nil")

type MongoRepository struct {
	db *mongo.Database
}

func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{
		db: db,
	}
}

func (mr *MongoRepository) CreateUser(user *model.User) error {
	if user == nil {
		logger.Logger.Error("user is nil")
		return fmt.Errorf("%w", ErrUserIsNil)
	}

	if !user.IsPasswordHashed {
		err := user.HashPassword()
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := mr.db.Collection("user").InsertOne(ctx, user)
	if err != nil {
		logger.Logger.Error("failed to create user", zap.Error(err))
		return fmt.Errorf("failed to create user: %w", err)
	}

	logger.Logger.Info("successfully created user", zap.String("username", user.Username))
	return nil
}
