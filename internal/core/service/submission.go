package service

import (
	"fmt"

	"github.com/SQL-Online-Judge/backend/internal/core/repository"
	"github.com/SQL-Online-Judge/backend/internal/model"
	"github.com/SQL-Online-Judge/backend/internal/pkg/logger"
	"github.com/SQL-Online-Judge/backend/internal/pkg/mq"
	"go.uber.org/zap"
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

	go func() {
		judgeRequest, err := ss.repo.GetJudgeRequest(submission)
		if err != nil {
			logger.Logger.Error("failed to get judge request", zap.Error(err))
			return
		}

		if !judgeRequest.Answer.IsReady {
			logger.Logger.Error("answer is not ready yet")
			return
		}

		judgeRequestJSON, err := judgeRequest.ToJSON()
		if err != nil {
			logger.Logger.Error("failed to marshal judge request", zap.Error(err))
			return
		}

		err = MQService.Enqueue(mq.QueueSubmission, judgeRequestJSON)
		if err != nil {
			logger.Logger.Error("failed to enqueue submission", zap.Error(err))
			return
		}

		ss.repo.UpdateSubmissionStatus(submissionID, model.JudgeStatusQueued)
	}()

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
