package initialize

import (
	"os"
	"strconv"

	"github.com/SQL-Online-Judge/backend/internal/core/repository"
	"github.com/SQL-Online-Judge/backend/internal/core/service"
	"github.com/SQL-Online-Judge/backend/internal/model"
	"github.com/SQL-Online-Judge/backend/internal/pkg/db/mongo"
	"github.com/SQL-Online-Judge/backend/internal/pkg/id"
	"github.com/SQL-Online-Judge/backend/internal/pkg/logger"
	"go.uber.org/zap"
)

func Initialize() {
	if hasInit, err := mongo.GetMongo().IsInitialized(); err != nil {
		logger.Logger.Fatal("failed to check if MongoDB has been initialized", zap.Error(err))
	} else if !hasInit {
		initMongo()
	} else {
		logger.Logger.Info("MongoDB has been initialized")
	}

	initRedis()
}

func initMongo() {
	logger.Logger.Info("initializing MongoDB...")
	createIndex()
	createAdmin()
	logger.Logger.Info("successfully initialized MongoDB")
}

func createIndex() {
	collectionIndexList := map[string][]map[string]string{
		"user": {
			{"field": "userID", "unique": "true"},
			{"field": "username", "unique": "true"},
			{"field": "deleted", "unique": "false"},
		},
		"class":      {{"field": "classID", "unique": "true"}},
		"problem":    {{"field": "problemID", "unique": "true"}},
		"answer":     {{"field": "answerID", "unique": "true"}},
		"task":       {{"field": "taskID", "unique": "true"}},
		"submission": {{"field": "submissionID", "unique": "true"}},
		"message":    {{"field": "messageID", "unique": "true"}},
		"messageBox": {{"field": "userID", "unique": "true"}},
	}

	for collection, indexList := range collectionIndexList {
		for _, indexMap := range indexList {
			unique, err := strconv.ParseBool(indexMap["unique"])
			if err != nil {
				logger.Logger.Fatal("failed to convert unique to bool", zap.Error(err))
			}

			err = mongo.GetMongo().CreateIndex(collection, indexMap["field"], unique)
			if err != nil {
				logger.Logger.Fatal("failed to create index in MongoDB", zap.Error(err))
			}
		}
	}
}

func createAdmin() {
	username := os.Getenv("ADMIN_USERNAME")
	password := os.Getenv("ADMIN_PASSWORD")

	admin := &model.User{
		UserID:   id.NewID(),
		Role:     "admin",
		Username: username,
		Password: password,
	}

	if !admin.IsValidUser() {
		logger.Logger.Fatal("invalid admin, please check the environment variables")
	}

	_, err := service.NewUserService(repository.NewMongoRepository(mongo.GetMongoDB())).CreateUser(admin.Username, admin.Password, admin.Role)
	if err != nil {
		logger.Logger.Fatal("failed to create admin", zap.Error(err))
	}

	logger.Logger.Info("successfully created admin")
}

func initRedis() {
	logger.Logger.Info("initializing Redis...")
	initRedisMQ()
	logger.Logger.Info("successfully initialized Redis")
}

func initRedisMQ() {
	queues := []string{"answer_generate", "answer_output", "submission", "judge_result"}

	for _, queue := range queues {
		exists, err := service.MQService.IsQueueExists(queue)
		if err != nil {
			logger.Logger.Fatal("failed to check if queue exists", zap.Error(err))
		}

		if exists {
			continue
		}

		err = service.MQService.CreateQueue(queue)
		if err != nil {
			logger.Logger.Fatal("failed to create queue", zap.Error(err))
		}

		logger.Logger.Info("successfully created queue", zap.String("queue", queue))
	}
}
