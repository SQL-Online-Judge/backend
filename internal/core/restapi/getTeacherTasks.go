package restapi

import (
	"net/http"
	"strconv"
	"time"

	"github.com/SQL-Online-Judge/backend/internal/pkg/logger"
	"go.uber.org/zap"
)

func getTeacherTasks(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	var resp getTasksResponse

	teacherID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		logger.Logger.Error("failed to get user id from context", zap.String("requestID", requestID))
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = &errorResponse{Code: http.StatusInternalServerError, Message: "failed to get user id from context"}
		w.Write(resp.toJSON())
		return
	}

	tasks, err := taskService.GetTeacherTasks(teacherID)
	if err != nil {
		logger.Logger.Error("failed to get teacher tasks", zap.Error(err), zap.String("requestID", requestID))
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = &errorResponse{Code: http.StatusInternalServerError, Message: "failed to get teacher tasks"}
		w.Write(resp.toJSON())
		return
	}

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
