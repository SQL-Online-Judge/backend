package service

import (
	"fmt"

	"github.com/SQL-Online-Judge/backend/internal/core/repository"
	"github.com/SQL-Online-Judge/backend/internal/model"
)

var (
	ErrSubmissionNotFound = fmt.Errorf("submission not found")
)

type SubmissionService struct {
	repo repository.SubmissionRepository
}

func NewSubmissionService(sr repository.SubmissionRepository) *SubmissionService {
	return &SubmissionService{
		repo: sr,
	}
}

func (ss *SubmissionService) CreateSubmission(submission *model.Submission) (int64, error) {
	submissionID, err := ss.repo.CreateSubmission(submission)
	if err != nil {
		return 0, fmt.Errorf("failed to create submission: %w", err)
	}

	return submissionID, nil
}

func (ss *SubmissionService) GetStudentSubmissions(studentID int64) ([]*model.SubmissionSummary, error) {
	submissions, err := ss.repo.FindSubmissionsByStudentID(studentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student submissions: %w", err)
	}

	return submissions, nil
}

func (ss *SubmissionService) isStudentSubmission(studentID, submissionID int64) bool {
	return ss.repo.IsStudentSubmission(studentID, submissionID)
}

func (ss *SubmissionService) GetStudentSubmittedSQL(studentID, submissionID int64) (string, error) {
	if !ss.isStudentSubmission(studentID, submissionID) {
		return "", fmt.Errorf("%w", ErrSubmissionNotFound)
	}

	sql, err := ss.repo.GetSubmittedSQL(submissionID)
	if err != nil {
		return "", fmt.Errorf("failed to get student submitted SQL: %w", err)
	}

	return sql, nil
}
