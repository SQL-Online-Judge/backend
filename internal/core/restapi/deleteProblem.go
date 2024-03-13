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

type deleteProblemResponse struct {
	ProblemID string         `json:"problemID,omitempty"`
	Error     *errorResponse `json:"error,omitempty"`
}

func (dpr *deleteProblemResponse) toJSON() []byte {
	res, err := json.Marshal(dpr)
	if err != nil {
		logger.Logger.Error("failed to marshal delete problem response", zap.Error(err))
		return nil
	}
	return res
}

func deleteProblem(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	var resp deleteProblemResponse

	sProblemID := chi.URLParam(r, "problemID")
	problemID, err := strconv.ParseInt(sProblemID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = &errorResponse{Code: http.StatusBadRequest, Message: "invalid problem id"}
		w.Write(resp.toJSON())
		return
	}

	teacherID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		logger.Logger.Error("failed to get teacher id from context", zap.String("requestID", requestID), zap.Any("teacherID", teacherID))
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = &errorResponse{Code: http.StatusInternalServerError, Message: "failed to get teacher id from context"}
		w.Write(resp.toJSON())
		return
	}

	err = problemService.DeleteProblem(teacherID, problemID)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		resp.ProblemID = sProblemID
		w.Write(resp.toJSON())
		return
	}

	switch {
	case errors.Is(err, service.ErrProblemNotFound):
		w.WriteHeader(http.StatusNotFound)
		resp.Error = &errorResponse{Code: http.StatusNotFound, Message: "problem not found"}
	case errors.Is(err, service.ErrNotProblemAuthor):
		w.WriteHeader(http.StatusForbidden)
		resp.Error = &errorResponse{Code: http.StatusForbidden, Message: "not the author of the problem"}
	default:
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = &errorResponse{Code: http.StatusInternalServerError, Message: "internal server error"}
	}
	w.Write(resp.toJSON())
}
