package id

import (
	"os"
	"strconv"

	"github.com/SQL-Online-Judge/backend/internal/pkg/logger"
	"github.com/bwmarrin/snowflake"
	"go.uber.org/zap"
)

var node *snowflake.Node

func init() {
	envNodeNumber := os.Getenv("SNOWFLAKE_NODE_NUMBER")
	if envNodeNumber == "" {
		logger.Logger.Fatal("SNOWFLAKE_NODE_NUMBER is not set")
	}

	nodeNumber, err := strconv.Atoi(envNodeNumber)
	if err != nil {
		logger.Logger.Fatal("failed to convert SNOWFLAKE_NODE_NUMBER to int", zap.Error(err))
	}

	node, err = snowflake.NewNode(int64(nodeNumber))
	if err != nil {
		logger.Logger.Fatal("failed to create snowflake node", zap.Error(err))
	}
	logger.Logger.Info("successfully initialized snowflake node", zap.Int("nodeNumber", nodeNumber))
}

func NewID() int64 {
	return node.Generate().Int64()
}
