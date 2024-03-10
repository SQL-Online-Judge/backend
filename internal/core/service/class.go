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
