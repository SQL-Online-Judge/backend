package restapi

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/SQL-Online-Judge/backend/internal/pkg/logger"
	"go.uber.org/zap"
)

type task struct {
	TaskID        string `json:"taskID"`
	TaskName      string `json:"taskName"`
	IsTimeLimited bool   `json:"isTimeLimited"`
	BeginTime     string `json:"beginTime"`
	EndTime       string `json:"endTime"`
}

type getTasksResponse struct {
	Tasks []task         `json:"tasks"`
	Error *errorResponse `json:"error,omitempty"`
}

func (gtr *getTasksResponse) toJSON() []byte {
	res, err := json.Marshal(gtr)
	if err != nil {
		logger.Logger.Error("failed to marshal get tasks response", zap.Error(err))
		return nil
	}
	return res
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	var resp getTasksResponse

	contains := r.URL.Query().Get("contains")
	tasks, err := taskService.GetTasks(contains)
	if err != nil {
		logger.Logger.Error("failed to get tasks", zap.Error(err), zap.String("requestID", requestID))
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = &errorResponse{Code: http.StatusInternalServerError, Message: "failed to get tasks"}
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
