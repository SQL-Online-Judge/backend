package service

import (
	"fmt"

	"github.com/SQL-Online-Judge/backend/internal/core/repository"
	"github.com/SQL-Online-Judge/backend/internal/model"
)

var (
	ErrAnswerNotFound     = fmt.Errorf("answer not found")
	ErrAnswerAlreadyExist = fmt.Errorf("answer already exist")
)

type AnswerService struct {
	repo repository.AnswerRepository
}

func NewAnswerService(ar repository.AnswerRepository) *AnswerService {
	return &AnswerService{
		repo: ar,
	}
}

func (as *AnswerService) isAnswerExist(problemID int64, dbName string) bool {
	return as.repo.IsAnswerExist(problemID, dbName)
}

func (as *AnswerService) CreateAnswer(teacherID int64, ps *ProblemService, answer *model.Answer) (int64, error) {
	problemID := answer.ProblemID

	if !ps.isProblemIDExist(problemID) {
		return 0, fmt.Errorf("%w", ErrProblemNotFound)
	}

	if ps.isProblemDeleted(problemID) {
		return 0, fmt.Errorf("%w", ErrProblemNotFound)
	}

	if !ps.checkProblemAuthor(teacherID, problemID) {
		return 0, fmt.Errorf("%w", ErrNotProblemAuthor)
	}

	if as.isAnswerExist(problemID, answer.DBName) {
		return 0, fmt.Errorf("%w", ErrAnswerAlreadyExist)
	}

	answerID, err := as.repo.CreateAnswer(answer)
	if err != nil {
		return 0, fmt.Errorf("failed to create answer: %w", err)
	}

	return answerID, nil
}
