package service

import (
	"fmt"

	"github.com/SQL-Online-Judge/backend/internal/core/repository"
	"github.com/SQL-Online-Judge/backend/internal/model"
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
