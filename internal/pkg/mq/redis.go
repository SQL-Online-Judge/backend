package mq

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const groupName = "sqloj"

var (
	ErrConsumerNameNotProvided = fmt.Errorf("consumer name is not provided")
	ErrConsumerNameNotString   = fmt.Errorf("consumer name is not a string")
	ErrNoMessageToDequeue      = fmt.Errorf("no message to dequeue")
	ErrDataFieldNotFound       = fmt.Errorf("data field not found")
	ErrDataFieldNotString      = fmt.Errorf("data field is not a string")
)

var Redis *RedisMQ

type RedisMQ struct {
	rdb *redis.Client
}

func NewRedisMQ(redis *redis.Client) *RedisMQ {
	return &RedisMQ{
		rdb: redis,
	}
}

func (r *RedisMQ) IsQueueExists(queueName string) (bool, error) {
	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()

	val, err := r.rdb.Exists(ctx, queueName).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check if queue exists: %w", err)
	}

	return val == 1, nil
}

func (r *RedisMQ) CreateQueue(queueName string) error {
	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()

	err := r.rdb.XGroupCreateMkStream(ctx, queueName, groupName, "0").Err()
	if err != nil {
		return fmt.Errorf("failed to create queue: %w", err)
	}

	return nil
}

func (r *RedisMQ) Enqueue(queueName, msg string) error {
	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()

	_, err := r.rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: queueName,
		Values: map[string]interface{}{"data": msg},
	}).Result()
	if err != nil {
		return fmt.Errorf("failed to enqueue message: %w", err)
	}

	return nil
}

func (r *RedisMQ) Dequeue(queueName string, args map[string]interface{}) (*Msg, error) {
	ctx := context.Background()

	iConsumerName, ok := args["consumerName"]
	if !ok {
		return nil, fmt.Errorf("%w", ErrConsumerNameNotProvided)
	}
	consumerName, ok := iConsumerName.(string)
	if !ok {
		return nil, fmt.Errorf("%w", ErrConsumerNameNotString)
	}

	res, err := r.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    groupName,
		Consumer: consumerName,
		Streams:  []string{queueName, ">"},
		Block:    -1,
		Count:    1,
	}).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to dequeue message: %w", err)
	}

	if len(res) == 0 || len(res[0].Messages) == 0 {
		return nil, fmt.Errorf("%w", ErrNoMessageToDequeue)
	}

	id := res[0].Messages[0].ID

	iData, ok := res[0].Messages[0].Values["data"]
	if !ok {
		return nil, fmt.Errorf("%w", ErrDataFieldNotFound)
	}

	data, ok := iData.(string)
	if !ok {
		return nil, fmt.Errorf("%w", ErrDataFieldNotString)
	}

	return &Msg{
		ID:   id,
		Data: data,
		Ack: func() error {
			return fmt.Errorf("failed to ack message: %w", r.rdb.XAck(ctx, queueName, groupName, id).Err())
		},
	}, nil
}
