package service

import (
	"github.com/SQL-Online-Judge/backend/internal/pkg/db/redis"
	"github.com/SQL-Online-Judge/backend/internal/pkg/mq"
)

var MQService *mq.Service = mq.NewService(mq.NewRedisMQ(redis.GetRedisDB()))
