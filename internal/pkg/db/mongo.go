package db

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/SQL-Online-Judge/backend/internal/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var Mongo *MongoDB
var ErrMongoURINotSet = errors.New("MONGO_URI is not set")
var ErrMongoClientNotConnected = errors.New("mongo client is not connected")

func init() {
	Mongo = &MongoDB{}
	err := Mongo.Connect(os.Getenv("MONGO_URI"))
	if err != nil {
		logger.Logger.Fatal("failed to connect to MongoDB during initialization", zap.Error(err))
	}
}

type MongoDB struct {
	client *mongo.Client
	db     *mongo.Database
}

func (m *MongoDB) Connect(uri string) error {
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

	u, err := url.Parse(uri)
	if err != nil {
		logger.Logger.Error("failed to parse MongoDB URI", zap.Error(err))
		return fmt.Errorf("failed to parse MongoDB URI: %w", err)
	}
	dbName := strings.TrimPrefix(u.Path, "/")
	m.client = client
	m.db = client.Database(dbName)
	logger.Logger.Info("successfully connected to MongoDB")

	return nil
}

func (m *MongoDB) Close() error {
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

func (m *MongoDB) Ping() error {
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

func (m *MongoDB) getCollectionNumber() (int, error) {
	if m.client == nil || m.db == nil {
		logger.Logger.Error(ErrMongoClientNotConnected.Error())
		return 0, fmt.Errorf("%w", ErrMongoClientNotConnected)
	}

	collections, err := m.db.ListCollectionNames(context.Background(), bson.D{})
	return len(collections), err
}

func (m *MongoDB) IsInitialized() (bool, error) {
	if m.client == nil || m.db == nil {
		logger.Logger.Error(ErrMongoClientNotConnected.Error())
		return false, fmt.Errorf("%w", ErrMongoClientNotConnected)
	}

	collectionNumber, err := m.getCollectionNumber()
	return collectionNumber > 0, err
}

func (m *MongoDB) GetDB() *mongo.Database {
	return m.db
}

func (m *MongoDB) GetCollection(name string) *mongo.Collection {
	return m.db.Collection(name)
}

func (m *MongoDB) CreateIndex(collectionName, fieldName string, unique bool) error {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: fieldName, Value: 1}},
		Options: options.Index().SetUnique(unique),
	}

	indexName, err := m.GetCollection(collectionName).Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		logger.Logger.Error("failed to create index",
			zap.String("collection", collectionName), zap.String("indexName", indexName), zap.Error(err))
		return fmt.Errorf("failed to create index: %w", err)
	}

	logger.Logger.Info("successfully created index",
		zap.String("collection", collectionName), zap.String("indexName", indexName))
	return nil
}

func GetMongo() *MongoDB {
	return Mongo
}
