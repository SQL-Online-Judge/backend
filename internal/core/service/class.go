package service

import (
	"fmt"

	"github.com/SQL-Online-Judge/backend/internal/core/repository"
	"github.com/SQL-Online-Judge/backend/internal/model"
)

var (
	ErrClassNotFound         = fmt.Errorf("class not found")
	ErrNotOfClassOwner       = fmt.Errorf("teacher is not the owner of the class")
	ErrStudentAlreadyInClass = fmt.Errorf("student is already in the class")
	ErrStudentNotInClass     = fmt.Errorf("student is not in the class")
	ErrTaskAlreadyInClass    = fmt.Errorf("task is already in the class")
	ErrTaskNotInClass        = fmt.Errorf("task is not in the class")
)

type ClassService struct {
	repo repository.ClassRepository
}

func NewClassService(cr repository.ClassRepository) *ClassService {
	return &ClassService{
		repo: cr,
	}
}

func (cs *ClassService) CreateClass(className string, teacherID int64) (int64, error) {
	classID, err := cs.repo.CreateClass(className, teacherID)
	if err != nil {
		return 0, fmt.Errorf("failed to create class: %w", err)
	}

	return classID, nil
}

func (cs *ClassService) isClassIDExist(classID int64) bool {
	return cs.repo.ExistByClassID(classID)
}

func (cs *ClassService) isClassDeleted(classID int64) bool {
	return cs.repo.IsClassDeleted(classID)
}

func (cs *ClassService) checkClassOwner(teacherID, classID int64) bool {
	return cs.repo.IsClassOwner(teacherID, classID)
}

func (cs *ClassService) DeleteClass(teacherID, classID int64) error {
	if !cs.isClassIDExist(classID) {
		return fmt.Errorf("%w", ErrClassNotFound)
	}

	if !cs.checkClassOwner(teacherID, classID) {
		return fmt.Errorf("%w", ErrNotOfClassOwner)
	}

	err := cs.repo.DeleteByClassID(classID)
	if err != nil {
		return fmt.Errorf("failed to delete class: %w", err)
	}

	return nil
}

func (cs *ClassService) UpdateClassName(teacherID, classID int64, className string) error {
	if !cs.isClassIDExist(classID) {
		return fmt.Errorf("%w", ErrClassNotFound)
	}

	if cs.isClassDeleted(classID) {
		return fmt.Errorf("%w", ErrClassNotFound)
	}

	if !cs.checkClassOwner(teacherID, classID) {
		return fmt.Errorf("%w", ErrNotOfClassOwner)
	}

	err := cs.repo.UpdateClassNameByClassID(classID, className)
	if err != nil {
		return fmt.Errorf("failed to update class name: %w", err)
	}

	return nil
}

func (cs *ClassService) GetClasses(teacherID int64) ([]*model.Class, error) {
	classes, err := cs.repo.FindClassesByTeacherID(teacherID)
	if err != nil {
		return nil, fmt.Errorf("failed to get classes: %w", err)
	}

	return classes, nil
}

func (cs *ClassService) isClassMember(classID, studentID int64) bool {
	return cs.repo.IsClassMember(classID, studentID)
}

func (cs *ClassService) AddStudentsToClass(us *UserService, teacherID, classID int64, studentIDs []int64) (map[int64]error, error) {
	if !cs.isClassIDExist(classID) {
		return nil, fmt.Errorf("%w", ErrClassNotFound)
	}

	if cs.isClassDeleted(classID) {
		return nil, fmt.Errorf("%w", ErrClassNotFound)
	}

	if !cs.checkClassOwner(teacherID, classID) {
		return nil, fmt.Errorf("%w", ErrNotOfClassOwner)
	}

	errs := make(map[int64]error)
	for _, studentID := range studentIDs {
		if err := us.isStudentExist(studentID); err != nil {
			errs[studentID] = err
			continue
		}

		if cs.isClassMember(classID, studentID) {
			errs[studentID] = fmt.Errorf("%w", ErrStudentAlreadyInClass)
			continue
		}

		err := cs.repo.AddStudentToClass(classID, studentID)
		if err != nil {
			errs[studentID] = fmt.Errorf("failed to add student to class: %w", err)
		} else {
			errs[studentID] = nil
		}
	}

	return errs, nil
}

func (cs *ClassService) RemoveStudentsFromClass(us *UserService, teacherID, classID int64, studentIDs []int64) (map[int64]error, error) {
	if !cs.isClassIDExist(classID) {
		return nil, fmt.Errorf("%w", ErrClassNotFound)
	}

	if cs.isClassDeleted(classID) {
		return nil, fmt.Errorf("%w", ErrClassNotFound)
	}

	if !cs.checkClassOwner(teacherID, classID) {
		return nil, fmt.Errorf("%w", ErrNotOfClassOwner)
	}

	errs := make(map[int64]error)
	for _, studentID := range studentIDs {
		if err := us.isStudentExist(studentID); err != nil {
			errs[studentID] = err
			continue
		}

		if !cs.isClassMember(classID, studentID) {
			errs[studentID] = fmt.Errorf("%w", ErrStudentNotInClass)
			continue
		}

		err := cs.repo.RemoveStudentFromClass(classID, studentID)
		if err != nil {
			errs[studentID] = fmt.Errorf("failed to remove student from class: %w", err)
		} else {
			errs[studentID] = nil
		}
	}

	return errs, nil
}

func (cs *ClassService) GetStudentsInClass(teacherID, classID int64) ([]*model.User, error) {
	if !cs.isClassIDExist(classID) {
		return nil, fmt.Errorf("%w", ErrClassNotFound)
	}

	if cs.isClassDeleted(classID) {
		return nil, fmt.Errorf("%w", ErrClassNotFound)
	}

	if !cs.checkClassOwner(teacherID, classID) {
		return nil, fmt.Errorf("%w", ErrNotOfClassOwner)
	}

	students, err := cs.repo.FindStudentsByClassID(classID)
	if err != nil {
		return nil, fmt.Errorf("failed to get students in class: %w", err)
	}

	return students, nil
}

func (cs *ClassService) isClassTask(classID, taskID int64) bool {
	return cs.repo.IsClassTask(classID, taskID)
}

func (cs *ClassService) AddTasks(ts *TaskService, teacherID, classID int64, taskIDs []int64) (map[int64]error, error) {
	if !cs.isClassIDExist(classID) {
		return nil, fmt.Errorf("%w", ErrClassNotFound)
	}

	if cs.isClassDeleted(classID) {
		return nil, fmt.Errorf("%w", ErrClassNotFound)
	}

	if !cs.checkClassOwner(teacherID, classID) {
		return nil, fmt.Errorf("%w", ErrNotOfClassOwner)
	}

	errs := make(map[int64]error)
	for _, taskID := range taskIDs {
		if !ts.isTaskIDExist(taskID) {
			errs[taskID] = fmt.Errorf("%w", ErrTaskNotFound)
			continue
		}

		if ts.isTaskDeleted(taskID) {
			errs[taskID] = fmt.Errorf("%w", ErrTaskNotFound)
			continue
		}

		if cs.isClassTask(classID, taskID) {
			errs[taskID] = fmt.Errorf("%w", ErrTaskAlreadyInClass)
			continue
		}

		err := cs.repo.AddTaskToClass(classID, taskID)
		if err != nil {
			errs[taskID] = fmt.Errorf("failed to add task to class: %w", err)
		} else {
			errs[taskID] = nil
		}
	}

	return errs, nil
}

func (cs *ClassService) RemoveTasks(ts *TaskService, teacherID, classID int64, taskIDs []int64) (map[int64]error, error) {
	if !cs.isClassIDExist(classID) {
		return nil, fmt.Errorf("%w", ErrClassNotFound)
	}

	if cs.isClassDeleted(classID) {
		return nil, fmt.Errorf("%w", ErrClassNotFound)
	}

	if !cs.checkClassOwner(teacherID, classID) {
		return nil, fmt.Errorf("%w", ErrNotOfClassOwner)
	}

	errs := make(map[int64]error)
	for _, taskID := range taskIDs {
		if !ts.isTaskIDExist(taskID) {
			errs[taskID] = fmt.Errorf("%w", ErrTaskNotFound)
			continue
		}

		if ts.isTaskDeleted(taskID) {
			errs[taskID] = fmt.Errorf("%w", ErrTaskNotFound)
			continue
		}

		if !cs.isClassTask(classID, taskID) {
			errs[taskID] = fmt.Errorf("%w", ErrTaskNotInClass)
			continue
		}

		err := cs.repo.RemoveTaskFromClass(classID, taskID)
		if err != nil {
			errs[taskID] = fmt.Errorf("failed to remove task from class: %w", err)
		} else {
			errs[taskID] = nil
		}
	}

	return errs, nil
}

func (cs *ClassService) GetTasksInClass(teacherID, classID int64) ([]*model.Task, error) {
	if !cs.isClassIDExist(classID) {
		return nil, fmt.Errorf("%w", ErrClassNotFound)
	}

	if cs.isClassDeleted(classID) {
		return nil, fmt.Errorf("%w", ErrClassNotFound)
	}

	if !cs.checkClassOwner(teacherID, classID) {
		return nil, fmt.Errorf("%w", ErrNotOfClassOwner)
	}

	tasks, err := cs.repo.GetTasksInClass(classID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks in class: %w", err)
	}

	return tasks, nil
}

func (cs *ClassService) GetClass(teacherID, classID int64) (*model.Class, []*model.User, []*model.Task, error) {
	if !cs.isClassIDExist(classID) {
		return nil, nil, nil, fmt.Errorf("%w", ErrClassNotFound)
	}

	if cs.isClassDeleted(classID) {
		return nil, nil, nil, fmt.Errorf("%w", ErrClassNotFound)
	}

	if !cs.checkClassOwner(teacherID, classID) {
		return nil, nil, nil, fmt.Errorf("%w", ErrNotOfClassOwner)
	}

	class, err := cs.repo.FindByClassID(classID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get class: %w", err)
	}

	students, err := cs.repo.FindStudentsByClassID(classID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get students in class: %w", err)
	}

	tasks, err := cs.repo.GetTasksInClass(classID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get tasks in class: %w", err)
	}

	return class, students, tasks, nil
}
