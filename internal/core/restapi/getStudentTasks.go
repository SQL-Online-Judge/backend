package restapi

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/SQL-Online-Judge/backend/internal/core/service"
	"github.com/SQL-Online-Judge/backend/internal/pkg/logger"
	"go.uber.org/zap"
)

func getStudentTasks(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	var resp getTasksResponse

	studentID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		logger.Logger.Error("failed to get user id from context", zap.String("requestID", requestID))
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = &errorResponse{Code: http.StatusInternalServerError, Message: "failed to get user id from context"}
		w.Write(resp.toJSON())
		return
	}

	tasks, err := taskService.GetStudentTasks(userService, studentID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			w.WriteHeader(http.StatusNotFound)
			resp.Error = &errorResponse{Code: http.StatusNotFound, Message: "student not found"}
		case errors.Is(err, service.ErrUserNotStudent):
			w.WriteHeader(http.StatusForbidden)
			resp.Error = &errorResponse{Code: http.StatusForbidden, Message: "user is not a student"}
		default:
			logger.Logger.Error("failed to get student tasks",
				zap.String("requestID", requestID),
				zap.Int64("studentID", studentID),
				zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			resp.Error = &errorResponse{Code: http.StatusInternalServerError, Message: "failed to get student tasks"}
		}
		w.Write(resp.toJSON())
		return
	}

	resp.Tasks = make([]task, 0, len(tasks))
	for _, t := range tasks {
		resp.Tasks = append(resp.Tasks, task{
			TaskID:        strconv.FormatInt(t.TaskID, 10),
			TaskName:      t.TaskName,
			IsTimeLimited: t.IsTimeLimited,
			BeginTime:     t.BeginTime.Format(time.RFC3339),
			EndTime:       t.EndTime.Format(time.RFC3339),
		})
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp.toJSON())
}
