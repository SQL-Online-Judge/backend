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

type updateClassTasksRequest struct {
	Tasks []string `json:"tasks"`
}

type updateClassTasksStatus struct {
	TaskID  string `json:"taskID"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type updateClassTasksResponse struct {
	ClassID string                   `json:"classID,omitempty"`
	Status  []updateClassTasksStatus `json:"status,omitempty"`
	Error   *errorResponse           `json:"error,omitempty"`
}

func (uctr *updateClassTasksResponse) toJSON() []byte {
	res, err := json.Marshal(uctr)
	if err != nil {
		logger.Logger.Error("failed to marshal update class tasks response", zap.Error(err))
		return nil
	}
	return res
}

func updateClassTasks(w http.ResponseWriter, r *http.Request, updateType string) {
	requestID := getRequestID(r)
	var resp updateClassTasksResponse

	sClassID := chi.URLParam(r, "classID")
	classID, err := strconv.ParseInt(sClassID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = &errorResponse{Code: http.StatusBadRequest, Message: "failed to parse class id"}
		w.Write(resp.toJSON())
		return
	}

	var req updateClassTasksRequest
	err = json.NewDecoder(r.Body).Decode(&req)
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

	taskIDs := make([]int64, 0, len(req.Tasks))
	for _, sTaskID := range req.Tasks {
		taskID, err := strconv.ParseInt(sTaskID, 10, 64)
		if err != nil {
			resp.Status = append(resp.Status, updateClassTasksStatus{TaskID: sTaskID, Code: http.StatusBadRequest, Message: "failed to parse task id"})
			continue
		}
		taskIDs = append(taskIDs, taskID)
	}

	var status map[int64]error
	switch updateType {
	case "add":
		status, err = classService.AddTasks(taskService, teacherID, classID, taskIDs)
	case "remove":
		status, err = classService.RemoveTasks(taskService, teacherID, classID, taskIDs)
	default:
		logger.Logger.Error("invalid update type", zap.String("requestID", requestID))
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = &errorResponse{Code: http.StatusInternalServerError, Message: "invalid update type"}
		w.Write(resp.toJSON())
		return
	}

	if err != nil {
		resp.Error = handleUpdateClassTasksError(w, err)
		resp.Status = nil
		w.Write(resp.toJSON())
		return
	}

	resp.ClassID = sClassID
	handleUpdateClassTasksStatus(status, &resp)
	w.WriteHeader(http.StatusOK)
	w.Write(resp.toJSON())
}

func handleUpdateClassTasksError(w http.ResponseWriter, err error) *errorResponse {
	var resp errorResponse
	switch {
	case errors.Is(err, service.ErrClassNotFound):
		w.WriteHeader(http.StatusNotFound)
		resp = errorResponse{Code: http.StatusNotFound, Message: "class not found"}
	case errors.Is(err, service.ErrNotOfClassOwner):
		w.WriteHeader(http.StatusForbidden)
		resp = errorResponse{Code: http.StatusForbidden, Message: "not class owner"}
	default:
		logger.Logger.Error("failed to update class tasks", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		resp = errorResponse{Code: http.StatusInternalServerError, Message: "failed to update class tasks"}
	}
	return &resp
}

func handleUpdateClassTasksStatus(status map[int64]error, resp *updateClassTasksResponse) {
	for taskID, err := range status {
		if err == nil {
			resp.Status = append(resp.Status, updateClassTasksStatus{
				TaskID:  strconv.FormatInt(taskID, 10),
				Code:    http.StatusOK,
				Message: "success",
			})
		} else {
			switch {
			case errors.Is(err, service.ErrTaskNotFound):
				resp.Status = append(resp.Status, updateClassTasksStatus{
					TaskID:  strconv.FormatInt(taskID, 10),
					Code:    http.StatusNotFound,
					Message: "task not found",
				})
			case errors.Is(err, service.ErrTaskAlreadyInClass):
				resp.Status = append(resp.Status, updateClassTasksStatus{
					TaskID:  strconv.FormatInt(taskID, 10),
					Code:    http.StatusConflict,
					Message: "task already in class",
				})
			case errors.Is(err, service.ErrTaskNotInClass):
				resp.Status = append(resp.Status, updateClassTasksStatus{
					TaskID:  strconv.FormatInt(taskID, 10),
					Code:    http.StatusConflict,
					Message: "task not in class",
				})
			default:
				logger.Logger.Error("failed to update class tasks", zap.Error(err))
				resp.Status = append(resp.Status, updateClassTasksStatus{
					TaskID:  strconv.FormatInt(taskID, 10),
					Code:    http.StatusInternalServerError,
					Message: "failed to update class tasks",
				})
			}
		}
	}
}

func addTasksToClass(w http.ResponseWriter, r *http.Request) {
	updateClassTasks(w, r, "add")
}

func removeTasksFromClass(w http.ResponseWriter, r *http.Request) {
	updateClassTasks(w, r, "remove")
}
