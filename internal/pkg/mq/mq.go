package mq

import (
	"fmt"

	"github.com/SQL-Online-Judge/backend/internal/pkg/logger"
	"go.uber.org/zap"
)

var (
	ErrQueueNameEmpty = fmt.Errorf("queue name is empty")
	ErrMessageEmpty   = fmt.Errorf("message is empty")
	ErrQueueNotFound  = fmt.Errorf("queue not found")
	ErrQueueExists    = fmt.Errorf("queue already exists")
)

type Msg struct {
	ID   string
	Data string
	Ack  func() error
}

type MQ interface {
	IsQueueExists(queueName string) (bool, error)
	CreateQueue(queueName string) error
	Enqueue(queueName, msg string) error
	Dequeue(queueName string, args map[string]interface{}) (*Msg, error)
}

const (
	QueueAnswerGenerate = "answer_generate"
	QueueAnswerOutput   = "answer_output"
	QueueSubmission     = "submission"
	QueueJudgeResult    = "judge_result"
)

type Service struct {
	mq MQ
}

func NewService(mq MQ) *Service {
	return &Service{
		mq: mq,
	}
}

func (ms *Service) IsQueueExists(queueName string) (bool, error) {
	if queueName == "" {
		logger.Logger.Error("queue name is empty")
		return false, fmt.Errorf("queue name is empty: %w", ErrQueueNameEmpty)
	}

	exists, err := ms.mq.IsQueueExists(queueName)
	if err != nil {
		logger.Logger.Error("failed to check if queue exists", zap.Error(err))
		return false, fmt.Errorf("failed to check if queue exists: %w", err)
	}

	return exists, nil
}

func (ms *Service) CreateQueue(queueName string) error {
	exists, err := ms.IsQueueExists(queueName)
	if err != nil {
		return err
	}

	if exists {
		logger.Logger.Error("queue already exists")
		return fmt.Errorf("queue already exists: %w", ErrQueueExists)
	}

	err = ms.mq.CreateQueue(queueName)
	if err != nil {
		logger.Logger.Error("failed to create queue", zap.Error(err))
		return fmt.Errorf("failed to create queue: %w", err)
	}

	return nil
}

func (ms *Service) Enqueue(queueName, msg string) error {
	err := ms.mq.Enqueue(queueName, msg)
	if err != nil {
		logger.Logger.Error("failed to enqueue message", zap.Error(err))
		return fmt.Errorf("failed to enqueue message: %w", err)
	}

	return nil
}

func (ms *Service) Dequeue(queueName string, args map[string]interface{}) (*Msg, error) {
	msg, err := ms.mq.Dequeue(queueName, args)
	if err != nil {
		logger.Logger.Error("failed to dequeue message", zap.Error(err))
		return nil, fmt.Errorf("failed to dequeue message: %w", err)
	}

	return msg, nil
}
