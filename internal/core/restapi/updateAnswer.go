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

type upadateAnswerRequest struct {
	PrepareSQL string `json:"prepareSQL"`
	AnswerSQL  string `json:"answerSQL"`
	JudgeSQL   string `json:"judgeSQL"`
}

func (uar *upadateAnswerRequest) toAnswer() *model.Answer {
	answer := &model.Answer{
		PrepareSQL: uar.PrepareSQL,
		AnswerSQL:  uar.AnswerSQL,
		JudgeSQL:   uar.JudgeSQL,
	}

	return answer
}

func (uar *upadateAnswerRequest) isValid() bool {
	answer := uar.toAnswer()
	return answer.IsValidPrepareSQL() && answer.IsValidAnswerSQL() && answer.IsValidJudgeSQL()
}

type upadateAnswerResponse struct {
	ProblemID string         `json:"problemID,omitempty"`
	AnswerID  string         `json:"answerID,omitempty"`
	Error     *errorResponse `json:"error,omitempty"`
}

func (uar *upadateAnswerResponse) toJSON() []byte {
	res, err := json.Marshal(uar)
	if err != nil {
		logger.Logger.Error("failed to marshal update answer response", zap.Error(err))
		return nil
	}

	return res
}

func updateAnswer(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	var resp upadateAnswerResponse

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

	var uar upadateAnswerRequest
	if err := json.NewDecoder(r.Body).Decode(&uar); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = &errorResponse{Code: http.StatusBadRequest, Message: "failed to decode update answer request"}
		w.Write(resp.toJSON())
		return
	}

	if !uar.isValid() {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = &errorResponse{Code: http.StatusBadRequest, Message: "invalid update answer request"}
		w.Write(resp.toJSON())
		return
	}

	answer := uar.toAnswer()
	answer.AnswerID = answerID
	answer.ProblemID = problemID

	err = answerService.UpdateAnswer(teacherID, problemService, answer)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		resp.ProblemID = sProblemID
		resp.AnswerID = sAnswerID
		w.Write(resp.toJSON())
		return
	}

	handleUpdateAnswerError(w, &resp, err)
}

func handleUpdateAnswerError(w http.ResponseWriter, resp *upadateAnswerResponse, err error) {
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
		logger.Logger.Error("failed to update answer", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = &errorResponse{Code: http.StatusInternalServerError, Message: "failed to update answer"}
	}

	w.Write(resp.toJSON())
}
