package main

import (
	"github.com/SQL-Online-Judge/backend/internal/core/restapi"
	"github.com/SQL-Online-Judge/backend/internal/pkg/db"
	"github.com/SQL-Online-Judge/backend/internal/pkg/logger"
)

func main() {
	logger.Logger.Info("Hello, SQL-Online-Judge!")

	defer db.Mongo.Close()
	defer db.Redis.Close()

	restapi.Serve()
}
