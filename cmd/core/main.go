package main

import (
	"github.com/SQL-Online-Judge/backend/internal/core/initialize"
	"github.com/SQL-Online-Judge/backend/internal/core/restapi"
	"github.com/SQL-Online-Judge/backend/internal/pkg/db/mongo"
	"github.com/SQL-Online-Judge/backend/internal/pkg/db/redis"
	"github.com/SQL-Online-Judge/backend/internal/pkg/logger"
)

func main() {
	defer mongo.GetMongo().Close()
	defer redis.GetRedis().Close()

	initialize.Initialize()
	logger.Logger.Info("Hello, SQL-Online-Judge!")

	restapi.Serve()
}
