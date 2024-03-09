package service

import (
	"fmt"

	"github.com/SQL-Online-Judge/backend/internal/core/repository"
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
