package restapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/SQL-Online-Judge/backend/internal/core/service"
	"github.com/SQL-Online-Judge/backend/internal/pkg/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type getStudentSubmittedSQLResponse struct {
	SubmissionID string         `json:"submissionID,omitempty"`
	SubmittedSQL string         `json:"submittedSQL,omitempty"`
	Error        *errorResponse `json:"error,omitempty"`
}

func (gsssr *getStudentSubmittedSQLResponse) toJSON() []byte {
	res, err := json.Marshal(gsssr)
	if err != nil {
		logger.Logger.Error("failed to marshal get student submitted sql response", zap.Error(err))
		return nil
	}
	return res
}

func getStudentSubmittedSQL(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	var resp getStudentSubmittedSQLResponse

	sSubmissionID := chi.URLParam(r, "submissionID")
	submissionID, err := strconv.ParseInt(sSubmissionID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = &errorResponse{Code: http.StatusBadRequest, Message: "invalid submission id"}
		w.Write(resp.toJSON())
		return
	}

	studentID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		logger.Logger.Error("failed to get user id from context", zap.String("requestID", requestID))
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = &errorResponse{Code: http.StatusInternalServerError, Message: "internal server error"}
		w.Write(resp.toJSON())
		return
	}

	sql, err := submissionService.GetStudentSubmittedSQL(studentID, submissionID)
	if err == nil {
		resp.SubmissionID = sSubmissionID
		resp.SubmittedSQL = sql

		w.WriteHeader(http.StatusOK)
		w.Write(resp.toJSON())
		return
	}

	switch {
	case errors.Is(err, service.ErrSubmissionNotFound):
		w.WriteHeader(http.StatusNotFound)
		resp.Error = &errorResponse{Code: http.StatusNotFound, Message: "submission not found"}
	default:
		logger.Logger.Error("failed to get student submitted sql", zap.Error(err), zap.String("requestID", requestID))
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = &errorResponse{Code: http.StatusInternalServerError, Message: "internal server error"}
	}

	w.Write(resp.toJSON())
}
