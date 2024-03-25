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

type studentTaskProblem struct {
	ProblemID string   `json:"problemID"`
	Title     string   `json:"title"`
	Tags      []string `json:"tags"`
	Score     string   `json:"score"`
}

type getStudentTaskProblemsResponse struct {
	TaskID   string                `json:"taskID,omitempty"`
	Problems []*studentTaskProblem `json:"problems,omitempty"`
	Error    *errorResponse        `json:"error,omitempty"`
}

func (gstpr *getStudentTaskProblemsResponse) toJSON() []byte {
	res, err := json.Marshal(gstpr)
	if err != nil {
		logger.Logger.Error("failed to marshal get student task problems response", zap.Error(err))
		return nil
	}
	return res
}

func getStudentTaskProblems(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	var resp getStudentTaskProblemsResponse

	sTaskID := chi.URLParam(r, "taskID")
	taskID, err := strconv.ParseInt(sTaskID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = &errorResponse{Code: http.StatusBadRequest, Message: "failed to parse task id"}
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

	taskProblems, problems, err := taskService.GetStudentTaskProblems(userService, studentID, taskID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			w.WriteHeader(http.StatusNotFound)
			resp.Error = &errorResponse{Code: http.StatusNotFound, Message: "student not found"}
		case errors.Is(err, service.ErrUserNotStudent):
			w.WriteHeader(http.StatusForbidden)
			resp.Error = &errorResponse{Code: http.StatusForbidden, Message: "user is not a student"}
		case errors.Is(err, service.ErrTaskNotFound):
			w.WriteHeader(http.StatusNotFound)
			resp.Error = &errorResponse{Code: http.StatusNotFound, Message: "task not found"}
		case errors.Is(err, service.ErrCannotAccessTask):
			w.WriteHeader(http.StatusForbidden)
			resp.Error = &errorResponse{Code: http.StatusForbidden, Message: "cannot access task"}
		default:
			logger.Logger.Error("failed to get student task problems",
				zap.String("requestID", requestID),
				zap.Int64("studentID", studentID),
				zap.Int64("taskID", taskID),
				zap.Error(err),
			)
			w.WriteHeader(http.StatusInternalServerError)
			resp.Error = &errorResponse{Code: http.StatusInternalServerError, Message: "failed to get student task problems"}
		}
		w.Write(resp.toJSON())
		return
	}

	taskProblemMap := make(map[int64]*model.TaskProblem)
	for _, taskProblem := range taskProblems {
		taskProblemMap[taskProblem.ProblemID] = taskProblem
	}

	resp.Problems = make([]*studentTaskProblem, 0, len(problems))
	for _, problem := range problems {
		taskProblem, ok := taskProblemMap[problem.ProblemID]
		if !ok {
			logger.Logger.Error("task problem not found",
				zap.String("requestID", requestID),
				zap.Int64("taskID", taskID),
				zap.Int64("problemID", problem.ProblemID),
			)
			w.WriteHeader(http.StatusInternalServerError)
			resp.Error = &errorResponse{Code: http.StatusInternalServerError, Message: "task problem not found"}
			w.Write(resp.toJSON())
			return
		}

		resp.Problems = append(resp.Problems, &studentTaskProblem{
			ProblemID: strconv.FormatInt(taskProblem.ProblemID, 10),
			Title:     problem.Title,
			Tags:      problem.Tags,
			Score:     strconv.FormatFloat(taskProblem.Score, 'f', -1, 64),
		})
	}

	resp.TaskID = sTaskID
	w.WriteHeader(http.StatusOK)
	w.Write(resp.toJSON())
}
