package service

import (
	"fmt"

	"github.com/SQL-Online-Judge/backend/internal/core/repository"
	"github.com/SQL-Online-Judge/backend/internal/model"
)

var (
	ErrProblemNotFound  = fmt.Errorf("problem not found")
	ErrNotProblemAuthor = fmt.Errorf("not the author of the problem")
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

func (ps *ProblemService) isProblemIDExist(problemID int64) bool {
	return ps.repo.ExistByProblemID(problemID)
}

func (ps *ProblemService) isProblemDeleted(problemID int64) bool {
	return ps.repo.IsProblemDeleted(problemID)
}

func (ps *ProblemService) checkProblemAuthor(teacherID, problemID int64) bool {
	return ps.repo.IsProblemAuthor(teacherID, problemID)
}

func (ps *ProblemService) DeleteProblem(teacherID, problemID int64) error {
	if !ps.isProblemIDExist(problemID) {
		return fmt.Errorf("%w", ErrProblemNotFound)
	}

	if ps.isProblemDeleted(problemID) {
		return fmt.Errorf("%w", ErrProblemNotFound)
	}

	if !ps.checkProblemAuthor(teacherID, problemID) {
		return fmt.Errorf("%w", ErrNotProblemAuthor)
	}

	err := ps.repo.DeleteByProblemID(problemID)
	if err != nil {
		return fmt.Errorf("failed to delete problem: %w", err)
	}

	return nil
}

func (ps *ProblemService) UpdateProblem(p *model.Problem) error {
	if !ps.isProblemIDExist(p.ProblemID) {
		return fmt.Errorf("%w", ErrProblemNotFound)
	}

	if ps.isProblemDeleted(p.ProblemID) {
		return fmt.Errorf("%w", ErrProblemNotFound)
	}

	if !ps.checkProblemAuthor(p.AuthorID, p.ProblemID) {
		return fmt.Errorf("%w", ErrNotProblemAuthor)
	}

	err := ps.repo.UpdateProblem(p)
	if err != nil {
		return fmt.Errorf("failed to update problem: %w", err)
	}

	return nil
}

func (ps *ProblemService) GetProblem(problemID int64) (*model.Problem, error) {
	if !ps.isProblemIDExist(problemID) {
		return nil, fmt.Errorf("%w", ErrProblemNotFound)
	}

	if ps.isProblemDeleted(problemID) {
		return nil, fmt.Errorf("%w", ErrProblemNotFound)
	}

	problem, err := ps.repo.FindByProblemID(problemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get problem: %w", err)
	}

	return problem, nil
}

func (ps *ProblemService) GetProblems(contains string) ([]*model.Problem, error) {
	problems, err := ps.repo.FindProblems(contains)

	if err != nil {
		return nil, fmt.Errorf("failed to get problems: %w", err)
	}

	return problems, nil
}

func (ps *ProblemService) GetTeacherProblems(teacherID int64) ([]*model.Problem, error) {
	problems, err := ps.repo.FindProblemsByAuthorID(teacherID)

	if err != nil {
		return nil, fmt.Errorf("failed to get teacher problems: %w", err)
	}

	return problems, nil
}
