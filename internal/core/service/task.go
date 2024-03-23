package service

import (
	"fmt"

	"github.com/SQL-Online-Judge/backend/internal/core/repository"
	"github.com/SQL-Online-Judge/backend/internal/model"
)

var (
	ErrTaskNotFound            = fmt.Errorf("task not found")
	ErrNotTaskAuthor           = fmt.Errorf("not the author of the task")
	ErrTaskProblemAlreadyExist = fmt.Errorf("task problem already exist")
	ErrTaskProblemNotFound     = fmt.Errorf("task problem not found")
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

func (ts *TaskService) isTaskProblem(taskID, problemID int64) bool {
	return ts.repo.IsTaskProblem(taskID, problemID)
}

func (ts *TaskService) AddTaskProblems(ps *ProblemService, teacherID, taskID int64, problems []*model.TaskProblem) (map[int64]error, error) {
	if !ts.isTaskIDExist(taskID) {
		return nil, fmt.Errorf("%w", ErrTaskNotFound)
	}

	if ts.isTaskDeleted(taskID) {
		return nil, fmt.Errorf("%w", ErrTaskNotFound)
	}

	if !ts.checkTaskAuthor(teacherID, taskID) {
		return nil, fmt.Errorf("%w", ErrNotTaskAuthor)
	}

	errs := make(map[int64]error)
	for _, problem := range problems {
		if !ps.isProblemIDExist(problem.ProblemID) {
			errs[problem.ProblemID] = fmt.Errorf("%w", ErrProblemNotFound)
			continue
		}

		if ps.isProblemDeleted(problem.ProblemID) {
			errs[problem.ProblemID] = fmt.Errorf("%w", ErrProblemNotFound)
			continue
		}

		if ts.isTaskProblem(taskID, problem.ProblemID) {
			errs[problem.ProblemID] = fmt.Errorf("%w", ErrTaskProblemAlreadyExist)
			continue
		}

		err := ts.repo.AddTaskProblem(taskID, problem)
		if err != nil {
			errs[problem.ProblemID] = fmt.Errorf("failed to add task problem: %w", err)
		} else {
			errs[problem.ProblemID] = nil
		}
	}

	return errs, nil
}

func (ts *TaskService) RemoveTaskProblems(ps *ProblemService, teacherID, taskID int64, problems []*model.TaskProblem) (map[int64]error, error) {
	if !ts.isTaskIDExist(taskID) {
		return nil, fmt.Errorf("%w", ErrTaskNotFound)
	}

	if ts.isTaskDeleted(taskID) {
		return nil, fmt.Errorf("%w", ErrTaskNotFound)
	}

	if !ts.checkTaskAuthor(teacherID, taskID) {
		return nil, fmt.Errorf("%w", ErrNotTaskAuthor)
	}

	errs := make(map[int64]error)
	for _, problem := range problems {
		if !ps.isProblemIDExist(problem.ProblemID) {
			errs[problem.ProblemID] = fmt.Errorf("%w", ErrProblemNotFound)
			continue
		}

		if ps.isProblemDeleted(problem.ProblemID) {
			errs[problem.ProblemID] = fmt.Errorf("%w", ErrProblemNotFound)
			continue
		}

		if !ts.isTaskProblem(taskID, problem.ProblemID) {
			errs[problem.ProblemID] = fmt.Errorf("%w", ErrTaskProblemNotFound)
			continue
		}

		err := ts.repo.RemoveTaskProblem(taskID, problem.ProblemID)
		if err != nil {
			errs[problem.ProblemID] = fmt.Errorf("failed to remove task problem: %w", err)
		} else {
			errs[problem.ProblemID] = nil
		}
	}

	return errs, nil
}

func (ts *TaskService) GetTask(taskID int64) (*model.Task, error) {
	if !ts.isTaskIDExist(taskID) {
		return nil, fmt.Errorf("%w", ErrTaskNotFound)
	}

	if ts.isTaskDeleted(taskID) {
		return nil, fmt.Errorf("%w", ErrTaskNotFound)
	}

	task, err := ts.repo.FindByTaskID(taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return task, nil
}

func (ts *TaskService) GetTasks(contains string) ([]*model.Task, error) {
	tasks, err := ts.repo.FindTasks(contains)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	return tasks, nil
}

func (ts *TaskService) GetTeacherTasks(teacherID int64) ([]*model.Task, error) {
	tasks, err := ts.repo.FindTasksByAuthorID(teacherID)
	if err != nil {
		return nil, fmt.Errorf("failed to get teacher tasks: %w", err)
	}

	return tasks, nil
}

func (ts *TaskService) GetStudentTasks(us *UserService, studentID int64) ([]*model.Task, error) {
	if err := us.isStudentExist(studentID); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	tasks, err := ts.repo.FindTasksByStudentID(studentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student tasks: %w", err)
	}

	return tasks, nil
}
