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

type getProblemResponse struct {
	ProblemID   string         `json:"problemID,omitempty"`
	Title       string         `json:"title,omitempty"`
	Tags        []string       `json:"tags,omitempty"`
	Content     string         `json:"content,omitempty"`
	TimeLimit   int32          `json:"timeLimit,omitempty"`
	MemoryLimit int32          `json:"memoryLimit,omitempty"`
	Error       *errorResponse `json:"error,omitempty"`
}

func (gpr *getProblemResponse) toJSON() []byte {
	res, err := json.Marshal(gpr)
	if err != nil {
		logger.Logger.Error("failed to marshal get problem response", zap.Error(err))
		return nil
	}
	return res
}

func getProblem(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	var resp getProblemResponse

	sProblemID := chi.URLParam(r, "problemID")
	problemID, err := strconv.ParseInt(sProblemID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = &errorResponse{Code: http.StatusBadRequest, Message: "invalid problem id"}
		w.Write(resp.toJSON())
		return
	}

	problem, err := problemService.GetProblem(problemID)
	if err == nil {
		resp.ProblemID = sProblemID
		resp.Title = problem.Title
		resp.Tags = problem.Tags
		resp.Content = problem.Content
		resp.TimeLimit = problem.TimeLimit
		resp.MemoryLimit = problem.MemoryLimit

		w.WriteHeader(http.StatusOK)
		w.Write(resp.toJSON())
		return
	}

	switch {
	case errors.Is(err, service.ErrProblemNotFound):
		w.WriteHeader(http.StatusNotFound)
		resp.Error = &errorResponse{Code: http.StatusNotFound, Message: "problem not found"}
	default:
		logger.Logger.Error("failed to get problem", zap.String("requestID", requestID), zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = &errorResponse{Code: http.StatusInternalServerError, Message: "internal server error"}
	}

	w.Write(resp.toJSON())
}
