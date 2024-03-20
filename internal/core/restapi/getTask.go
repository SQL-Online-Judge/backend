package restapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/SQL-Online-Judge/backend/internal/core/service"
	"github.com/SQL-Online-Judge/backend/internal/pkg/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type getTaskResponse struct {
	TaskID        string         `json:"taskID"`
	TaskName      string         `json:"taskName"`
	Problems      []taskProblem  `json:"problems"`
	IsTimeLimited bool           `json:"isTimeLimited"`
	BeginTime     string         `json:"beginTime"`
	EndTime       string         `json:"endTime"`
	Error         *errorResponse `json:"error,omitempty"`
}

func (gtr *getTaskResponse) toJSON() []byte {
	res, err := json.Marshal(gtr)
	if err != nil {
		logger.Logger.Error("failed to marshal get task response", zap.Error(err))
		return nil
	}
	return res
}

func getTask(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	var resp getTaskResponse

	sTaskID := chi.URLParam(r, "taskID")
	taskID, err := strconv.ParseInt(sTaskID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = &errorResponse{Code: http.StatusBadRequest, Message: "failed to parse task id"}
		w.Write(resp.toJSON())
		return
	}

	task, err := taskService.GetTask(taskID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrTaskNotFound):
			w.WriteHeader(http.StatusNotFound)
			resp.Error = &errorResponse{Code: http.StatusNotFound, Message: "task not found"}
		default:
			logger.Logger.Error("failed to get task", zap.Error(err), zap.String("requestID", requestID))
			w.WriteHeader(http.StatusInternalServerError)
			resp.Error = &errorResponse{Code: http.StatusInternalServerError, Message: "failed to get task"}
		}
		w.Write(resp.toJSON())
		return
	}

	resp.TaskID = sTaskID
	resp.TaskName = task.TaskName
	resp.Problems = make([]taskProblem, 0, len(task.Problems))
	for _, problem := range task.Problems {
		resp.Problems = append(resp.Problems, taskProblem{
			ProblemID: strconv.FormatInt(problem.ProblemID, 10),
			Score:     strconv.FormatFloat(problem.Score, 'f', -1, 64),
		})
	}
	resp.IsTimeLimited = task.IsTimeLimited
	resp.BeginTime = task.BeginTime.Format(time.RFC3339)
	resp.EndTime = task.EndTime.Format(time.RFC3339)
	w.WriteHeader(http.StatusOK)
	w.Write(resp.toJSON())
}
