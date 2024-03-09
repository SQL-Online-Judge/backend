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

type updateClassNameRequest struct {
	ClassName string `json:"className"`
}

type updateClassNameResponse struct {
	ClassID      string         `json:"classID,omitempty"`
	ClassNewName string         `json:"classNewName,omitempty"`
	Error        *errorResponse `json:"error,omitempty"`
}

func (ucr *updateClassNameResponse) toJSON() []byte {
	res, err := json.Marshal(ucr)
	if err != nil {
		logger.Logger.Error("failed to marshal update class name response", zap.Error(err))
		return nil
	}
	return res
}

func updateClassName(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	var resp updateClassNameResponse

	sClassID := chi.URLParam(r, "classID")
	classID, err := strconv.ParseInt(sClassID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = &errorResponse{Code: http.StatusBadRequest, Message: "invalid class id"}
		w.Write(resp.toJSON())
		return
	}

	var req updateClassNameRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = &errorResponse{Code: http.StatusBadRequest, Message: "failed to decode request body"}
		w.Write(resp.toJSON())
		return
	}

	class := model.Class{ClassID: classID, ClassName: req.ClassName}
	if !class.IsValidClassName() {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = &errorResponse{Code: http.StatusBadRequest, Message: "invalid class name"}
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

	err = classService.UpdateClassName(teacherID, classID, req.ClassName)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		resp.ClassID = sClassID
		resp.ClassNewName = req.ClassName
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
		logger.Logger.Error("failed to update class name", zap.String("requestID", requestID), zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = &errorResponse{Code: http.StatusInternalServerError, Message: "failed to update class name"}
	}

	w.Write(resp.toJSON())
}
