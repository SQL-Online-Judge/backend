package initialize

import (
	"strconv"

	"github.com/SQL-Online-Judge/backend/internal/core/repository"
	"github.com/SQL-Online-Judge/backend/internal/core/service"
	"github.com/SQL-Online-Judge/backend/internal/model"
	"github.com/SQL-Online-Judge/backend/internal/pkg/db"
	"github.com/SQL-Online-Judge/backend/internal/pkg/id"
	"github.com/SQL-Online-Judge/backend/internal/pkg/logger"
	"github.com/SQL-Online-Judge/backend/internal/pkg/random"
	"go.uber.org/zap"
)

func Initialize() {
	if hasInit, err := db.GetMongo().IsInitialized(); err != nil {
		logger.Logger.Fatal("failed to check if MongoDB has been initialized", zap.Error(err))
	} else if !hasInit {
		logger.Logger.Info("initializing MongoDB...")
		initMongo()
		logger.Logger.Info("successfully initialized MongoDB")
	} else {
		logger.Logger.Info("MongoDB has been initialized")
	}
}

func initMongo() {
	createIndex()
	createAdmin()
}

func createIndex() {
	collectionIndex := map[string]map[string]string{
		"user":             {"field": "userID", "unique": "true"},
		"class":            {"field": "classID", "unique": "true"},
		"problem":          {"field": "problemID", "unique": "true"},
		"answer":           {"field": "answerID", "unique": "true"},
		"problemSet":       {"field": "problemSetID", "unique": "true"},
		"classProblemSets": {"field": "classProblemSetID", "unique": "true"},
		"task":             {"field": "taskID", "unique": "true"},
		"submission":       {"field": "submissionID", "unique": "true"},
		"message":          {"field": "messageID", "unique": "true"},
		"messageBox":       {"field": "userID", "unique": "true"},
	}

	for collection, indexMap := range collectionIndex {
		unique, err := strconv.ParseBool(indexMap["unique"])
		if err != nil {
			logger.Logger.Fatal("failed to convert unique to bool", zap.Error(err))
		}

		err = db.GetMongo().CreateIndex(collection, indexMap["field"], unique)
		if err != nil {
			logger.Logger.Fatal("failed to create index in MongoDB", zap.Error(err))
		}
	}
}

func createAdmin() {
	username, err := random.NewRandomString(8)
	if err != nil {
		logger.Logger.Fatal("failed to generate random string", zap.Error(err))
	}
	password, err := random.NewRandomString(16)
	if err != nil {
		logger.Logger.Fatal("failed to generate random string", zap.Error(err))
	}

	admin := &model.User{
		UserID:   id.NewID(),
		Role:     "admin",
		Username: username,
		Password: password,
	}

	err = service.NewUserService(repository.NewMongoRepository(db.GetMongoDB())).CreateUser(admin)
	if err != nil {
		logger.Logger.Fatal("failed to create admin", zap.Error(err))
	}

	logger.Logger.Info("successfully created admin", zap.String("username", username), zap.String("password", password))
}
