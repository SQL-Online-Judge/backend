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

func (ts *TaskService) isTaskIDExist(taskID int64) bool {
	return ts.repo.ExistByTaskID(taskID)
}

func (ts *TaskService) isTaskDeleted(taskID int64) bool {
	return ts.repo.IsTaskDeleted(taskID)
}

func (ts *TaskService) checkTaskAuthor(teacherID, taskID int64) bool {
	return ts.repo.IsTaskAuthor(teacherID, taskID)
}

func (ts *TaskService) DeleteTask(teacherID, taskID int64) error {
	if !ts.isTaskIDExist(taskID) {
		return fmt.Errorf("%w", ErrTaskNotFound)
	}

	if ts.isTaskDeleted(taskID) {
		return fmt.Errorf("%w", ErrTaskNotFound)
	}

	if !ts.checkTaskAuthor(teacherID, taskID) {
		return fmt.Errorf("%w", ErrNotTaskAuthor)
	}

	err := ts.repo.DeleteByTaskID(taskID)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}

func (ts *TaskService) UpdateTask(task *model.Task) error {
	if !ts.isTaskIDExist(task.TaskID) {
		return fmt.Errorf("%w", ErrTaskNotFound)
	}

	if ts.isTaskDeleted(task.TaskID) {
		return fmt.Errorf("%w", ErrTaskNotFound)
	}

	if !ts.checkTaskAuthor(task.AuthorID, task.TaskID) {
		return fmt.Errorf("%w", ErrNotTaskAuthor)
	}

	err := ts.repo.UpdateTask(task)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}
