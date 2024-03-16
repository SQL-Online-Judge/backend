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

type deleteTaskResponse struct {
	TaskID string         `json:"taskID,omitempty"`
	Error  *errorResponse `json:"error,omitempty"`
}

func (dtr *deleteTaskResponse) toJSON() []byte {
	res, err := json.Marshal(dtr)
	if err != nil {
		logger.Logger.Error("failed to marshal delete task response", zap.Error(err))
		return nil
	}

	return res
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	var resp deleteTaskResponse

	sTaskID := chi.URLParam(r, "taskID")
	taskID, err := strconv.ParseInt(sTaskID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = &errorResponse{Code: http.StatusBadRequest, Message: "failed to parse taskID"}
		w.Write(resp.toJSON())
		return
	}

	teacherID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		logger.Logger.Error("failed to get user id from context", zap.String("requestID", requestID))
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = &errorResponse{Code: http.StatusInternalServerError, Message: "failed to get user id from context"}
		w.Write(resp.toJSON())
		return
	}

	err = taskService.DeleteTask(teacherID, taskID)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		resp.TaskID = sTaskID
		w.Write(resp.toJSON())
		return
	}

	switch {
	case errors.Is(err, service.ErrTaskNotFound):
		w.WriteHeader(http.StatusNotFound)
		resp.Error = &errorResponse{Code: http.StatusNotFound, Message: "task not found"}
	case errors.Is(err, service.ErrNotTaskAuthor):
		w.WriteHeader(http.StatusForbidden)
		resp.Error = &errorResponse{Code: http.StatusForbidden, Message: "not the author of the task"}
	default:
		logger.Logger.Error("failed to delete task", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = &errorResponse{Code: http.StatusInternalServerError, Message: "failed to delete task"}
	}

	w.Write(resp.toJSON())
}
