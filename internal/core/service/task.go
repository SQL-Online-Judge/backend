package service

import (
	"fmt"

	"github.com/SQL-Online-Judge/backend/internal/core/repository"
	"github.com/SQL-Online-Judge/backend/internal/model"
)

var (
	ErrTaskNotFound  = fmt.Errorf("task not found")
	ErrNotTaskAuthor = fmt.Errorf("not the author of the task")
)

type TaskService struct {
	repo repository.TaskRepository
}

func NewTaskService(tr repository.TaskRepository) *TaskService {
	return &TaskService{
		repo: tr,
	}
}

func (ts *TaskService) CreateTask(task *model.Task) (int64, error) {
	taskID, err := ts.repo.CreateTask(task)
	if err != nil {
		return 0, fmt.Errorf("failed to create task: %w", err)
	}

	return taskID, nil
}
