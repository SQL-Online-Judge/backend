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

type getStudentResponse struct {
	UserID   string         `json:"userID,omitempty"`
	Username string         `json:"username,omitempty"`
	Error    *errorResponse `json:"error,omitempty"`
}

func (gsr *getStudentResponse) toJSON() []byte {
	res, err := json.Marshal(gsr)
	if err != nil {
		logger.Logger.Error("failed to marshal get student response", zap.Error(err))
		return nil
	}
	return res
}

func getStudent(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	var resp getStudentResponse

	sUserID := chi.URLParam(r, "userID")
	userID, err := strconv.ParseInt(sUserID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = &errorResponse{Code: http.StatusBadRequest, Message: "invalid user id"}
		w.Write(resp.toJSON())
		return
	}

	student, err := userService.GetStudent(userID)
	if err == nil {
		resp.UserID = strconv.FormatInt(student.UserID, 10)
		resp.Username = student.Username
		w.WriteHeader(http.StatusOK)
		w.Write(resp.toJSON())
		return
	}

	switch {
	case errors.Is(err, service.ErrUserNotFound):
		w.WriteHeader(http.StatusNotFound)
		resp.Error = &errorResponse{Code: http.StatusNotFound, Message: "user not found"}
	case errors.Is(err, service.ErrUserNotStudent):
		w.WriteHeader(http.StatusForbidden)
		resp.Error = &errorResponse{Code: http.StatusForbidden, Message: "user is not student"}
	default:
		logger.Logger.Error("failed to get student",
			zap.String("requestID", requestID),
			zap.String("userID", sUserID),
			zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = &errorResponse{Code: http.StatusInternalServerError, Message: "internal server error"}
	}

	w.Write(resp.toJSON())
}
