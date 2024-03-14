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

type createAnswerRequest struct {
	DBName     string `json:"dbName"`
	PrepareSQL string `json:"prepareSQL"`
	AnswerSQL  string `json:"answerSQL"`
	JudgeSQL   string `json:"judgeSQL"`
}

type createAnswerResponse struct {
	AnswerID string         `json:"answerID,omitempty"`
	Error    *errorResponse `json:"error,omitempty"`
}

func (car *createAnswerResponse) toJSON() []byte {
	res, err := json.Marshal(car)
	if err != nil {
		logger.Logger.Error("failed to marshal create answer response", zap.Error(err))
		return nil
	}
	return res
}

func createAnswer(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	var resp createAnswerResponse

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

	var req createAnswerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = &errorResponse{Code: http.StatusBadRequest, Message: "invalid request body"}
		w.Write(resp.toJSON())
		return
	}

	answer := model.NewAnswer(&model.Answer{
		ProblemID:  problemID,
		DBName:     req.DBName,
		PrepareSQL: req.PrepareSQL,
		AnswerSQL:  req.AnswerSQL,
		JudgeSQL:   req.JudgeSQL,
	})

	if !answer.IsValidAnswer() {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = &errorResponse{Code: http.StatusBadRequest, Message: "invalid answer"}
		w.Write(resp.toJSON())
		return
	}

	answerID, err := answerService.CreateAnswer(teacherID, problemService, answer)
	if err == nil {
		resp.AnswerID = strconv.FormatInt(answerID, 10)
		w.WriteHeader(http.StatusOK)
		w.Write(resp.toJSON())
		return
	}

	switch {
	case errors.Is(err, service.ErrProblemNotFound):
		w.WriteHeader(http.StatusNotFound)
		resp.Error = &errorResponse{Code: http.StatusNotFound, Message: "problem not found"}
	case errors.Is(err, service.ErrNotProblemAuthor):
		w.WriteHeader(http.StatusForbidden)
		resp.Error = &errorResponse{Code: http.StatusForbidden, Message: "not problem author"}
	case errors.Is(err, service.ErrAnswerAlreadyExist):
		w.WriteHeader(http.StatusConflict)
		resp.Error = &errorResponse{Code: http.StatusConflict, Message: "answer already exist"}
	default:
		logger.Logger.Error("failed to create answer", zap.String("requestID", requestID), zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = &errorResponse{Code: http.StatusInternalServerError, Message: "internal server error"}
	}

	w.Write(resp.toJSON())
}
