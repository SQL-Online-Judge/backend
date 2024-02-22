package db

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/SQL-Online-Judge/backend/internal/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var Mongo *mongoDB

var ErrMongoURINotSet = errors.New("MONGO_URI is not set")

var ErrMongoClientNotConnected = errors.New("mongo client is not connected")

func init() {
	Mongo = &mongoDB{}

	err := Mongo.Connect(os.Getenv("MONGO_URI"))
	if err != nil {
		logger.Logger.Fatal("failed to connect to MongoDB during initialization", zap.Error(err))
	}
}

type mongoDB struct {
	client *mongo.Client
	db     *mongo.Database
}

func (m *mongoDB) Connect(uri string) error {
	var err error

	if uri == "" {
		logger.Logger.Error(ErrMongoURINotSet.Error())

		return fmt.Errorf("%w", ErrMongoURINotSet)
	}

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		logger.Logger.Error("failed to connect to MongoDB", zap.Error(err))

		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		logger.Logger.Error("failed to ping MongoDB", zap.Error(err))

		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	m.client = client
	m.db = client.Database("")

	logger.Logger.Info("successfully connected to MongoDB")

	return nil
}

func (m *mongoDB) Close() error {
	if m.client == nil {
		logger.Logger.Warn("mongo client is not connected")

		return nil
	}

	err := m.client.Disconnect(context.Background())
	if err != nil {
		logger.Logger.Error("failed to disconnect from MongoDB", zap.Error(err))

		return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
	}

	m.client = nil

	return nil
}

func (m *mongoDB) Ping() error {
	if m.client == nil {
		logger.Logger.Error(ErrMongoClientNotConnected.Error())

		return fmt.Errorf("%w", ErrMongoClientNotConnected)
	}

	err := m.client.Ping(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return nil
}
