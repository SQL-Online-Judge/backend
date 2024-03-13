package service

import (
	"fmt"

	"github.com/SQL-Online-Judge/backend/internal/core/repository"
	"github.com/SQL-Online-Judge/backend/internal/model"
)

var (
	ErrProblemNotFound = fmt.Errorf("problem not found")
)

type ProblemService struct {
	repo repository.ProblemRepository
}

func NewProblemService(pr repository.ProblemRepository) *ProblemService {
	return &ProblemService{
		repo: pr,
	}
}

func (ps *ProblemService) CreateProblem(p *model.Problem) (int64, error) {
	problemID, err := ps.repo.CreateProblem(p)
	if err != nil {
		return 0, fmt.Errorf("failed to create problem: %w", err)
	}

	return problemID, nil
}
