package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/SQL-Online-Judge/backend/internal/model"
	"github.com/SQL-Online-Judge/backend/internal/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
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

func (mr *MongoRepository) getCollection() *mongo.Collection {
	return mr.db.Collection("user")
}

func (mr *MongoRepository) ExistByUserID(userID int64) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.D{{Key: "userID", Value: userID}}
	count, err := mr.getCollection().CountDocuments(ctx, filter)
	if err != nil {
		logger.Logger.Error("failed to count documents", zap.Error(err))
		return false
	}

	return count > 0
}

func (mr *MongoRepository) ExistByUsername(username string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.D{{Key: "username", Value: username}}
	count, err := mr.getCollection().CountDocuments(ctx, filter)
	if err != nil {
		logger.Logger.Error("failed to count documents", zap.Error(err))
		return false
	}

	return count > 0
}

func (mr *MongoRepository) Create(username, password, role string) (int64, error) {
	user := model.NewUser(username, password, role)
	hashedPassword, err := user.GetHashedPassword()
	if err != nil {
		return 0, fmt.Errorf("failed to get hashed password: %w", err)
	}
	user.Password = hashedPassword

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = mr.getCollection().InsertOne(ctx, user)
	if err != nil {
		logger.Logger.Error("failed to create user", zap.String("username", user.Username), zap.Error(err))
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return user.UserID, nil
}

func (mr *MongoRepository) FindByUserID(userID int64) (*model.User, error) {
	var user model.User
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.D{
		{Key: "userID", Value: userID},
		{Key: "deleted", Value: false},
	}
	err := mr.getCollection().FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by user id: %w", err)
	}

	return &user, nil
}

func (mr *MongoRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.D{
		{Key: "username", Value: username},
		{Key: "deleted", Value: false},
	}
	err := mr.getCollection().FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by username: %w", err)
	}

	return &user, nil
}

func (mr *MongoRepository) DeleteByUserID(userID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.D{{Key: "userID", Value: userID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "deleted", Value: true}}}}
	_, err := mr.getCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Logger.Error("failed to delete user", zap.Int64("userID", userID), zap.Error(err))
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
