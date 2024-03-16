package restapi

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/SQL-Online-Judge/backend/internal/model"
	"github.com/SQL-Online-Judge/backend/internal/pkg/logger"
	"go.uber.org/zap"
)

type updateTaskRequest struct {
	TaskName      string    `json:"taskName"`
	IsTimeLimited bool      `json:"isTimeLimited"`
	BeginTime     time.Time `json:"beginTime"`
	EndTime       time.Time `json:"endTime"`
}

func (ctr *updateTaskRequest) toTask() *model.Task {
	return &model.Task{
		TaskName:      ctr.TaskName,
		IsTimeLimited: ctr.IsTimeLimited,
		BeginTime:     ctr.BeginTime,
		EndTime:       ctr.EndTime,
	}
}

type updateTaskResponse struct {
	TaskID string         `json:"taskID,omitempty"`
	Error  *errorResponse `json:"error,omitempty"`
}

func (ctr *updateTaskResponse) toJSON() []byte {
	res, err := json.Marshal(ctr)
	if err != nil {
		logger.Logger.Error("failed to marshal create task response", zap.Error(err))
		return nil
	}

	return res
}

func createTask(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	var resp updateTaskResponse

	var req updateTaskRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = &errorResponse{Code: http.StatusBadRequest, Message: "failed to decode request body"}
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

	task := req.toTask()
	task.AuthorID = teacherID
	task = model.NewTask(task)
	if !task.IsValidTask() {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = &errorResponse{Code: http.StatusBadRequest, Message: "invalid task"}
		w.Write(resp.toJSON())
		return
	}

	taskID, err := taskService.CreateTask(task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = &errorResponse{Code: http.StatusInternalServerError, Message: "failed to create task"}
		w.Write(resp.toJSON())
		return
	}

	w.WriteHeader(http.StatusOK)
	resp.TaskID = strconv.FormatInt(taskID, 10)
	w.Write(resp.toJSON())
}
