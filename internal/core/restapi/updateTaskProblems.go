package restapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/SQL-Online-Judge/backend/internal/core/service"
	"github.com/SQL-Online-Judge/backend/internal/model"
	"github.com/SQL-Online-Judge/backend/internal/pkg/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type taskProblem struct {
	ProblemID string `json:"problemID"`
	Score     string `json:"score,omitempty"`
}

type updateTaskProblemRequest struct {
	Problems []taskProblem `json:"problems"`
}

type updateTaskProblemStatus struct {
	ProblemID string `json:"problemID"`
	Code      int    `json:"code"`
	Message   string `json:"message"`
}

type updateTaskProblemResponse struct {
	TaskID string                    `json:"taskID,omitempty"`
	Status []updateTaskProblemStatus `json:"status,omitempty"`
	Error  *errorResponse            `json:"error,omitempty"`
}

func (utpr *updateTaskProblemResponse) toJSON() []byte {
	res, err := json.Marshal(utpr)
	if err != nil {
		logger.Logger.Error("failed to marshal update task problem response", zap.Error(err))
		return nil
	}
	return res
}

func updateTaskProblem(w http.ResponseWriter, r *http.Request, updateType string) {
	requestID := getRequestID(r)
	var resp updateTaskProblemResponse

	sTaskID := chi.URLParam(r, "taskID")
	taskID, err := strconv.ParseInt(sTaskID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = &errorResponse{Code: http.StatusBadRequest, Message: "failed to parse task id"}
		w.Write(resp.toJSON())
		return
	}

	var req updateTaskProblemRequest
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

	problems := make([]*model.TaskProblem, 0, len(req.Problems))
	for _, problem := range req.Problems {
		problemID, err := strconv.ParseInt(problem.ProblemID, 10, 64)
		if err != nil {
			resp.Status = append(resp.Status, updateTaskProblemStatus{
				ProblemID: problem.ProblemID,
				Code:      http.StatusBadRequest,
				Message:   "failed to parse problem id",
			})
			continue
		}
		score, err := strconv.ParseFloat(problem.Score, 64)
		if updateType == "add" && err != nil {
			resp.Status = append(resp.Status, updateTaskProblemStatus{
				ProblemID: problem.ProblemID,
				Code:      http.StatusBadRequest,
				Message:   "failed to parse score",
			})
			continue
		}
		taskProblem := model.TaskProblem{
			ProblemID: problemID,
			Score:     score,
		}
		if updateType == "add" && !taskProblem.IsValidScore() {
			resp.Status = append(resp.Status, updateTaskProblemStatus{
				ProblemID: problem.ProblemID,
				Code:      http.StatusBadRequest,
				Message:   "invalid score",
			})
			continue
		}
		problems = append(problems, &taskProblem)
	}

	var status map[int64]error
	switch updateType {
	case "add":
		status, err = taskService.AddTaskProblems(problemService, teacherID, taskID, problems)
	case "remove":
		status, err = taskService.RemoveTaskProblems(problemService, teacherID, taskID, problems)
	default:
		logger.Logger.Error("invalid update type", zap.String("requestID", requestID), zap.String("updateType", updateType))
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = &errorResponse{Code: http.StatusInternalServerError, Message: "invalid update type"}
		w.Write(resp.toJSON())
		return
	}

	if err != nil {
		resp.Error = handleUpdateTaskProblemError(w, err)
		resp.Status = nil
		w.Write(resp.toJSON())
		return
	}

	resp.TaskID = sTaskID
	handleUpdateTaskProblemStatus(status, &resp)
	w.WriteHeader(http.StatusOK)
	w.Write(resp.toJSON())
}

func handleUpdateTaskProblemError(w http.ResponseWriter, err error) *errorResponse {
	var resp *errorResponse
	switch {
	case errors.Is(err, service.ErrTaskNotFound):
		w.WriteHeader(http.StatusNotFound)
		resp = &errorResponse{Code: http.StatusNotFound, Message: "task not found"}
	case errors.Is(err, service.ErrNotTaskAuthor):
		w.WriteHeader(http.StatusForbidden)
		resp = &errorResponse{Code: http.StatusForbidden, Message: "not the author of the task"}
	default:
		logger.Logger.Error("failed to update task problem", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		resp = &errorResponse{Code: http.StatusInternalServerError, Message: "failed to update task problem"}
	}
	return resp
}

func handleUpdateTaskProblemStatus(status map[int64]error, resp *updateTaskProblemResponse) {
	for problemID, err := range status {
		if err == nil {
			resp.Status = append(resp.Status, updateTaskProblemStatus{
				ProblemID: strconv.FormatInt(problemID, 10),
				Code:      http.StatusOK,
				Message:   "success",
			})
			continue
		}
		switch {
		case errors.Is(err, service.ErrProblemNotFound) || errors.Is(err, service.ErrTaskProblemNotFound):
			resp.Status = append(resp.Status, updateTaskProblemStatus{
				ProblemID: strconv.FormatInt(problemID, 10),
				Code:      http.StatusNotFound,
				Message:   "problem not found",
			})
		case errors.Is(err, service.ErrTaskProblemAlreadyExist):
			resp.Status = append(resp.Status, updateTaskProblemStatus{
				ProblemID: strconv.FormatInt(problemID, 10),
				Code:      http.StatusConflict,
				Message:   "task problem already exist",
			})
		default:
			logger.Logger.Error("failed to update task problem", zap.Error(err))
			resp.Status = append(resp.Status, updateTaskProblemStatus{
				ProblemID: strconv.FormatInt(problemID, 10),
				Code:      http.StatusInternalServerError,
				Message:   "failed to update task problem",
			})
		}
	}
}
