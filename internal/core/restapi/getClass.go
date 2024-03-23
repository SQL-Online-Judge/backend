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

type getClassResponse struct {
	ClassID   string         `json:"classID,omitempty"`
	ClassName string         `json:"className,omitempty"`
	Students  []student      `json:"students,omitempty"`
	Tasks     []task         `json:"tasks,omitempty"`
	Error     *errorResponse `json:"error,omitempty"`
}

func (gcr *getClassResponse) toJSON() []byte {
	res, err := json.Marshal(gcr)
	if err != nil {
		logger.Logger.Error("failed to marshal get class response", zap.Error(err))
		return nil
	}
	return res
}

func getClass(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	var resp getClassResponse

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

	class, students, tasks, err := classService.GetClass(teacherID, classID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrClassNotFound):
			w.WriteHeader(http.StatusNotFound)
			resp.Error = &errorResponse{Code: http.StatusNotFound, Message: "class not found"}
		case errors.Is(err, service.ErrNotOfClassOwner):
			w.WriteHeader(http.StatusForbidden)
			resp.Error = &errorResponse{Code: http.StatusForbidden, Message: "not the owner of the class"}
		default:
			logger.Logger.Error("failed to get class", zap.String("requestID", requestID), zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			resp.Error = &errorResponse{Code: http.StatusInternalServerError, Message: "failed to get class"}
		}
		w.Write(resp.toJSON())
		return
	}

	resp.ClassID = strconv.FormatInt(class.ClassID, 10)
	resp.ClassName = class.ClassName
	for _, s := range students {
		resp.Students = append(resp.Students, student{
			UserID:   strconv.FormatInt(s.UserID, 10),
			Username: s.Username,
		})
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
