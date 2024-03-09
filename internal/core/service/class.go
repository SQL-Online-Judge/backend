package service

import (
	"fmt"

	"github.com/SQL-Online-Judge/backend/internal/core/repository"
)

var (
	ErrClassNotFound   = fmt.Errorf("class not found")
	ErrNotOfClassOwner = fmt.Errorf("teacher is not the owner of the class")
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
