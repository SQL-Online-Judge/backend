package db

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/SQL-Online-Judge/backend/internal/pkg/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var Redis *RedisDB
var ErrRedisURLNotSet = errors.New("REDIS_URL is not set")
var ErrRedisClientNotConnected = errors.New("redis client is not connected")

func init() {
	Redis = &RedisDB{}

	err := Redis.Connect(os.Getenv("REDIS_URL"))
	if err != nil {
		logger.Logger.Fatal("failed to connect to Redis during initialization", zap.Error(err))
	}
}

type RedisDB struct {
	client *redis.Client
}

func (r *RedisDB) Connect(url string) error {
	if url == "" {
		logger.Logger.Error(ErrRedisURLNotSet.Error())
		return fmt.Errorf("%w", ErrRedisURLNotSet)
	}

	clientOptions, err := redis.ParseURL(url)
	if err != nil {
		logger.Logger.Error("failed to parse redis url", zap.Error(err))
		return fmt.Errorf("failed to parse redis url: %w", err)
	}

	client := redis.NewClient(clientOptions)

	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		logger.Logger.Error("failed to ping Redis", zap.Error(err))
		return fmt.Errorf("failed to ping Redis: %w", err)
	}

	r.client = client
	logger.Logger.Info("successfully connected to Redis")
	return nil
}

func (r *RedisDB) Close() error {
	if r.client == nil {
		logger.Logger.Warn(ErrRedisClientNotConnected.Error())
		return nil
	}

	err := r.client.Close()
	if err != nil {
		logger.Logger.Error("failed to disconnect from Redis", zap.Error(err))
		return fmt.Errorf("failed to disconnect from Redis: %w", err)
	}

	r.client = nil
	return nil
}

func (r *RedisDB) Ping() error {
	if r.client == nil {
		logger.Logger.Error(ErrRedisClientNotConnected.Error())
		return fmt.Errorf("%w", ErrRedisClientNotConnected)
	}

	_, err := r.client.Ping(context.Background()).Result()
	if err != nil {
		logger.Logger.Error("failed to ping Redis", zap.Error(err))
		return fmt.Errorf("failed to ping Redis: %w", err)
	}

	logger.Logger.Info("successfully pinged Redis")
	return nil
}

func GetRedis() *RedisDB {
	return Redis
}

func GetRedisDB() *redis.Client {
	return Redis.client
}
