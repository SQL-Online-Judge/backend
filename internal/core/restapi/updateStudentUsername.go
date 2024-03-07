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

type updateStudentUsernameRequest struct {
	Username string `json:"username"`
}

type updateStudentUsernameResponse struct {
	UserID   string         `json:"userID,omitempty"`
	Username string         `json:"username,omitempty"`
	Error    *errorResponse `json:"error,omitempty"`
}

func (usr *updateStudentUsernameResponse) toJSON() []byte {
	res, err := json.Marshal(usr)
	if err != nil {
		logger.Logger.Error("failed to marshal update student username response", zap.Error(err))
		return nil
	}
	return res
}

func updateStudentUsername(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	var req updateStudentUsernameRequest
	var resp updateStudentUsernameResponse

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = &errorResponse{Code: http.StatusBadRequest, Message: "invalid request"}
		w.Write(resp.toJSON())
		return
	}

	sUserID := chi.URLParam(r, "userID")
	userID, err := strconv.ParseInt(sUserID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = &errorResponse{Code: http.StatusBadRequest, Message: "invalid user id"}
		w.Write(resp.toJSON())
		return
	}

	student := &model.User{Username: req.Username}
	if !student.IsValidUsername() {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = &errorResponse{Code: http.StatusBadRequest, Message: "invalid username"}
		w.Write(resp.toJSON())
		return
	}

	err = userService.UpdateStudentUsername(userID, req.Username)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		resp.UserID = sUserID
		resp.Username = req.Username
		w.Write(resp.toJSON())
		return
	}

	switch {
	case errors.Is(err, service.ErrUserNotFound):
		w.WriteHeader(http.StatusNotFound)
		resp.Error = &errorResponse{Code: http.StatusNotFound, Message: "user not found"}
	case errors.Is(err, service.ErrUserConflict):
		w.WriteHeader(http.StatusConflict)
		resp.Error = &errorResponse{Code: http.StatusConflict, Message: "username already exists"}
	case errors.Is(err, service.ErrUserNotStudent):
		w.WriteHeader(http.StatusForbidden)
		resp.Error = &errorResponse{Code: http.StatusForbidden, Message: "user is not a student"}
	default:
		logger.Logger.Error("failed to update student username",
			zap.String("requestID", requestID),
			zap.String("userID", sUserID),
			zap.String("username", req.Username),
			zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = &errorResponse{Code: http.StatusInternalServerError, Message: "failed to update student username"}
	}

	w.Write(resp.toJSON())
}
