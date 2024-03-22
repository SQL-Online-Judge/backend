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

type getTaksInClassResponse struct {
	ClassID string         `json:"classID,omitempty"`
	Tasks   []task         `json:"tasks,omitempty"`
	Error   *errorResponse `json:"error,omitempty"`
}

func (gticr *getTaksInClassResponse) toJSON() []byte {
	res, err := json.Marshal(gticr)
	if err != nil {
		logger.Logger.Error("failed to marshal get tasks in class response", zap.Error(err))
		return nil
	}
	return res
}

func getTasksInClass(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	var resp getTaksInClassResponse

	sClassID := chi.URLParam(r, "classID")
	classID, err := strconv.ParseInt(sClassID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = &errorResponse{Code: http.StatusBadRequest, Message: "failed to parse class id"}
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

	tasks, err := classService.GetTasksInClass(teacherID, classID)
	if err == nil {
		resp.ClassID = strconv.FormatInt(classID, 10)
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
		return
	}

	switch {
	case errors.Is(err, service.ErrClassNotFound):
		w.WriteHeader(http.StatusNotFound)
		resp.Error = &errorResponse{Code: http.StatusNotFound, Message: "class not found"}
	case errors.Is(err, service.ErrNotOfClassOwner):
		w.WriteHeader(http.StatusForbidden)
		resp.Error = &errorResponse{Code: http.StatusForbidden, Message: "not of class owner"}
	default:
		logger.Logger.Error("failed to get tasks in class", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = &errorResponse{Code: http.StatusInternalServerError, Message: "failed to get tasks in class"}
	}

	w.Write(resp.toJSON())
}
