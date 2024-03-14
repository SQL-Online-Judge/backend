package service

import (
	"fmt"

	"github.com/SQL-Online-Judge/backend/internal/core/repository"
	"github.com/SQL-Online-Judge/backend/internal/model"
)

var (
	ErrAnswerNotFound     = fmt.Errorf("answer not found")
	ErrAnswerAlreadyExist = fmt.Errorf("answer already exist")
	ErrNotAnswerOfProblem = fmt.Errorf("not the answer of the problem")
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

func (as *AnswerService) isAnswerIDExist(answerID int64) bool {
	return as.repo.ExistByAnswerID(answerID)
}

func (as *AnswerService) isAnswerDeleted(answerID int64) bool {
	return as.repo.IsAnswerDeleted(answerID)
}

func (as *AnswerService) isAnswerOfProblem(problemID, answerID int64) bool {
	return as.repo.IsAnswerOfProblem(problemID, answerID)
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

func (as *AnswerService) checkUpdateAnswer(teacherID int64, ps *ProblemService, problemID, answerID int64) error {
	if !ps.isProblemIDExist(problemID) {
		return fmt.Errorf("%w", ErrProblemNotFound)
	}

	if ps.isProblemDeleted(problemID) {
		return fmt.Errorf("%w", ErrProblemNotFound)
	}

	if !ps.checkProblemAuthor(teacherID, problemID) {
		return fmt.Errorf("%w", ErrNotProblemAuthor)
	}

	if !as.isAnswerIDExist(answerID) {
		return fmt.Errorf("%w", ErrAnswerNotFound)
	}

	if as.isAnswerDeleted(answerID) {
		return fmt.Errorf("%w", ErrAnswerNotFound)
	}

	if !as.isAnswerOfProblem(problemID, answerID) {
		return fmt.Errorf("%w", ErrNotAnswerOfProblem)
	}

	return nil
}

func (as *AnswerService) DeleteAnswer(teacherID int64, ps *ProblemService, problemID, answerID int64) error {
	err := as.checkUpdateAnswer(teacherID, ps, problemID, answerID)
	if err != nil {
		return err
	}

	err = as.repo.DeleteByAnswerID(answerID)
	if err != nil {
		return fmt.Errorf("failed to delete answer: %w", err)
	}

	return nil
}

func (as *AnswerService) UpdateAnswer(teacherID int64, ps *ProblemService, answer *model.Answer) error {
	problemID := answer.ProblemID
	answerID := answer.AnswerID

	err := as.checkUpdateAnswer(teacherID, ps, problemID, answerID)
	if err != nil {
		return err
	}

	err = as.repo.UpdateAnswer(answer)
	if err != nil {
		return fmt.Errorf("failed to update answer: %w", err)
	}

	return nil
}
