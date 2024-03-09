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

type deleteClassResponse struct {
	ClassID string         `json:"classID,omitempty"`
	Error   *errorResponse `json:"error,omitempty"`
}

func (dcr *deleteClassResponse) toJSON() []byte {
	res, err := json.Marshal(dcr)
	if err != nil {
		logger.Logger.Error("failed to marshal delete class response", zap.Error(err))
		return nil
	}
	return res
}

func deleteClass(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	var resp deleteClassResponse

	sClassID := chi.URLParam(r, "classID")
	classID, err := strconv.ParseInt(sClassID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = &errorResponse{Code: http.StatusBadRequest, Message: "invalid class id"}
		w.Write(resp.toJSON())
		return
	}

	teacherID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		logger.Logger.Error("failed to get teacher id from context", zap.String("requestID", requestID), zap.Any("teacherID", teacherID))
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = &errorResponse{Code: http.StatusInternalServerError, Message: "failed to get teacher id from context"}
		w.Write(resp.toJSON())
		return
	}

	err = classService.DeleteClass(teacherID, classID)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		resp.ClassID = sClassID
		w.Write(resp.toJSON())
		return
	}

	switch {
	case errors.Is(err, service.ErrClassNotFound):
		w.WriteHeader(http.StatusNotFound)
		resp.Error = &errorResponse{Code: http.StatusNotFound, Message: "class not found"}
	case errors.Is(err, service.ErrNotOfClassOwner):
		w.WriteHeader(http.StatusForbidden)
		resp.Error = &errorResponse{Code: http.StatusForbidden, Message: "teacher is not the owner of the class"}
	default:
		logger.Logger.Error("failed to delete class", zap.String("requestID", requestID), zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = &errorResponse{Code: http.StatusInternalServerError, Message: "failed to delete class"}
	}

	w.Write(resp.toJSON())
}
