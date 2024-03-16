package restapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/SQL-Online-Judge/backend/internal/core/service"
	"github.com/SQL-Online-Judge/backend/internal/model"
	"github.com/SQL-Online-Judge/backend/internal/pkg/logger"
	"go.uber.org/zap"
)

type createTeacherRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (ctr *createTeacherRequest) toUser() *model.User {
	return &model.User{
		Username: ctr.Username,
		Password: ctr.Password,
		Role:     "teacher",
	}
}

type createTeacherResponse struct {
	UserID string         `json:"userID,omitempty"`
	Error  *errorResponse `json:"error,omitempty"`
}

func (ctr *createTeacherResponse) toJSON() []byte {
	res, err := json.Marshal(ctr)
	if err != nil {
		logger.Logger.Error("failed to marshal create teacher response", zap.Error(err))
		return nil
	}
	return res
}

func createTeacher(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	var req createTeacherRequest
	var teacher *model.User

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp := createTeacherResponse{Error: &errorResponse{Code: http.StatusBadRequest, Message: "invalid request"}}
		w.Write(resp.toJSON())
		return
	}

	teacher = req.toUser()
	if !teacher.IsValidUser() {
		w.WriteHeader(http.StatusBadRequest)
		resp := createTeacherResponse{Error: &errorResponse{Code: http.StatusBadRequest, Message: "invalid username or password"}}
		w.Write(resp.toJSON())
		return
	}

	teacher.UserID, err = userService.CreateUser(teacher.Username, teacher.Password, teacher.Role)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserConflict):
			w.WriteHeader(http.StatusConflict)
			resp := createTeacherResponse{Error: &errorResponse{Code: http.StatusConflict, Message: "username already exists"}}
			w.Write(resp.toJSON())
			return
		default:
			logger.Logger.Error("failed to create teacher", zap.String("requestID", requestID), zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			resp := createTeacherResponse{Error: &errorResponse{Code: http.StatusInternalServerError, Message: "failed to create teacher"}}
			w.Write(resp.toJSON())
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	resp := createTeacherResponse{UserID: strconv.FormatInt(teacher.UserID, 10)}
	w.Write(resp.toJSON())
}
