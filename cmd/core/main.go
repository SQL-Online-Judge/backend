package main

import (
	"github.com/SQL-Online-Judge/backend/internal/core/initialize"
	"github.com/SQL-Online-Judge/backend/internal/core/restapi"
	"github.com/SQL-Online-Judge/backend/internal/pkg/db"
	"github.com/SQL-Online-Judge/backend/internal/pkg/logger"
)

func main() {
	defer db.GetMongo().Close()
	defer db.GetRedis().Close()
	initialize.Initialize()
	logger.Logger.Info("Hello, SQL-Online-Judge!")

	restapi.Serve()
}
