package restapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/SQL-Online-Judge/backend/internal/core/service"
	"github.com/SQL-Online-Judge/backend/internal/model"
	"github.com/SQL-Online-Judge/backend/internal/pkg/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type createStudentSubmissionRequest struct {
	DBName       string `json:"dbName"`
	SubmittedSQL string `json:"submittedSQL"`
}

type createStudentSubmissionResponse struct {
	SubmissionID string         `json:"submissionID,omitempty"`
	Error        *errorResponse `json:"error,omitempty"`
}

func (cssr *createStudentSubmissionResponse) toJSON() []byte {
	res, err := json.Marshal(cssr)
	if err != nil {
		logger.Logger.Error("failed to marshal create student submission response", zap.Error(err))
		return nil
	}
	return res
}

func createStudentSubmission(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	var resp createStudentSubmissionResponse

	var req createStudentSubmissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = &errorResponse{Code: http.StatusBadRequest, Message: "failed to decode request body"}
		w.Write(resp.toJSON())
		return
	}

	sTaskID := chi.URLParam(r, "taskID")
	taskID, err := strconv.ParseInt(sTaskID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = &errorResponse{Code: http.StatusBadRequest, Message: "failed to parse task id"}
		w.Write(resp.toJSON())
		return
	}

	sProblemID := chi.URLParam(r, "problemID")
	problemID, err := strconv.ParseInt(sProblemID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = &errorResponse{Code: http.StatusBadRequest, Message: "invalid problem id"}
		w.Write(resp.toJSON())
		return
	}

	studentID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		logger.Logger.Error("failed to get user id from context", zap.String("requestID", requestID))
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = &errorResponse{Code: http.StatusInternalServerError, Message: "failed to get user id from context"}
		w.Write(resp.toJSON())
		return
	}

	submission := &model.Submission{
		SubmitterID:  studentID,
		TaskID:       taskID,
		ProblemID:    problemID,
		DBName:       req.DBName,
		SubmittedSQL: req.SubmittedSQL,
	}
	submission = model.NewSubmission(submission)

	if !submission.IsValidSubmission() {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = &errorResponse{Code: http.StatusBadRequest, Message: "invalid submission"}
		w.Write(resp.toJSON())
		return
	}

	submissionID, err := taskService.CreateStudentSubmission(userService, problemService, submissionService, submission)
	if err == nil {
		resp.SubmissionID = strconv.FormatInt(submissionID, 10)
		w.WriteHeader(http.StatusOK)
		w.Write(resp.toJSON())
		return
	}

	switch {
	case errors.Is(err, service.ErrTaskNotFound):
		w.WriteHeader(http.StatusNotFound)
		resp.Error = &errorResponse{Code: http.StatusNotFound, Message: "task not found"}
	case errors.Is(err, service.ErrCannotAccessTask):
		w.WriteHeader(http.StatusForbidden)
		resp.Error = &errorResponse{Code: http.StatusForbidden, Message: "cannot access task"}
	case errors.Is(err, service.ErrTaskProblemNotFound):
		w.WriteHeader(http.StatusNotFound)
		resp.Error = &errorResponse{Code: http.StatusNotFound, Message: "task problem not found"}
	case errors.Is(err, service.ErrProblemNotFound):
		w.WriteHeader(http.StatusNotFound)
		resp.Error = &errorResponse{Code: http.StatusNotFound, Message: "problem not found"}
	case errors.Is(err, service.ErrNotInSubmitTime):
		w.WriteHeader(http.StatusForbidden)
		resp.Error = &errorResponse{Code: http.StatusForbidden, Message: "not in submit time"}
	default:
		logger.Logger.Error("failed to create student submission", zap.Error(err), zap.String("requestID", requestID))
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = &errorResponse{Code: http.StatusInternalServerError, Message: "failed to create student submission"}
	}
	w.Write(resp.toJSON())
}
