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

type deleteAnswerResponse struct {
	ProblemID string         `json:"problemID,omitempty"`
	AnswerID  string         `json:"answerID,omitempty"`
	Error     *errorResponse `json:"error,omitempty"`
}

func (dar *deleteAnswerResponse) toJSON() []byte {
	res, err := json.Marshal(dar)
	if err != nil {
		logger.Logger.Error("failed to marshal delete answer response", zap.Error(err))
		return nil
	}
	return res
}

func deleteAnswer(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	var resp deleteAnswerResponse

	teacherID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		logger.Logger.Error("failed to get teacher id from context", zap.String("requestID", requestID), zap.Any("teacherID", teacherID))
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = &errorResponse{Code: http.StatusInternalServerError, Message: "failed to get teacher id from context"}
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

	sAnswerID := chi.URLParam(r, "answerID")
	answerID, err := strconv.ParseInt(sAnswerID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = &errorResponse{Code: http.StatusBadRequest, Message: "invalid answer id"}
		w.Write(resp.toJSON())
		return
	}

	err = answerService.DeleteAnswer(teacherID, problemService, problemID, answerID)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		resp.ProblemID = sProblemID
		resp.AnswerID = sAnswerID
		w.Write(resp.toJSON())
		return
	}

	switch {
	case errors.Is(err, service.ErrProblemNotFound):
		w.WriteHeader(http.StatusNotFound)
		resp.Error = &errorResponse{Code: http.StatusNotFound, Message: "problem not found"}
	case errors.Is(err, service.ErrNotProblemAuthor):
		w.WriteHeader(http.StatusForbidden)
		resp.Error = &errorResponse{Code: http.StatusForbidden, Message: "not the problem author"}
	case errors.Is(err, service.ErrAnswerNotFound):
		w.WriteHeader(http.StatusNotFound)
		resp.Error = &errorResponse{Code: http.StatusNotFound, Message: "answer not found"}
	case errors.Is(err, service.ErrNotAnswerOfProblem):
		w.WriteHeader(http.StatusForbidden)
		resp.Error = &errorResponse{Code: http.StatusForbidden, Message: "not the answer of the problem"}
	default:
		logger.Logger.Error("failed to delete answer", zap.String("requestID", requestID), zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = &errorResponse{Code: http.StatusInternalServerError, Message: "failed to delete answer"}
	}

	w.Write(resp.toJSON())
}
