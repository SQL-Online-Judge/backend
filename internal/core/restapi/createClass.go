package restapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/SQL-Online-Judge/backend/internal/model"
	"github.com/SQL-Online-Judge/backend/internal/pkg/logger"
	"go.uber.org/zap"
)

type createClassRequest struct {
	ClassName string `json:"className"`
}

type createClassResponse struct {
	ClassID   string         `json:"classID,omitempty"`
	ClassName string         `json:"className,omitempty"`
	Error     *errorResponse `json:"error,omitempty"`
}

func (ccr *createClassResponse) toJSON() []byte {
	res, err := json.Marshal(ccr)
	if err != nil {
		logger.Logger.Error("failed to marshal create class response", zap.Error(err))
		return nil
	}
	return res
}

func createClass(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	var req createClassRequest
	var resp createClassResponse

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = &errorResponse{Code: http.StatusBadRequest, Message: "invalid request"}
		w.Write(resp.toJSON())
		return
	}

	class := &model.Class{ClassName: req.ClassName}
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

	classID, err := classService.CreateClass(req.ClassName, teacherID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = &errorResponse{Code: http.StatusInternalServerError, Message: "failed to create class"}
		w.Write(resp.toJSON())
		return
	}

	resp.ClassID = strconv.FormatInt(classID, 10)
	resp.ClassName = req.ClassName
	w.Write(resp.toJSON())
}
