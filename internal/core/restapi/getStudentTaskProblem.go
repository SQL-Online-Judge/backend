package restapi

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/SQL-Online-Judge/backend/internal/core/service"
	"github.com/SQL-Online-Judge/backend/internal/pkg/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func getStudentTaskProblem(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	var resp getProblemResponse

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

	problem, err := taskService.GetStudentTaskProblem(userService, problemService, studentID, taskID, problemID)
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
	default:
		logger.Logger.Error("failed to get problem", zap.Error(err), zap.String("requestID", requestID))
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = &errorResponse{Code: http.StatusInternalServerError, Message: "internal server error"}
	}
	w.Write(resp.toJSON())
}
