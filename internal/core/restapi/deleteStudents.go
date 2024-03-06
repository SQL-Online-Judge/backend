package restapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/SQL-Online-Judge/backend/internal/core/service"
	"github.com/SQL-Online-Judge/backend/internal/pkg/logger"
	"go.uber.org/zap"
)

type deleteStudentsRequest struct {
	UserIDs []string `json:"userIDs"`
}

type deleteStudentStatus struct {
	UserID  string `json:"userID"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type deleteStudentsResponse struct {
	Status []deleteStudentStatus `json:"status,omitempty"`
	Error  *errorResponse        `json:"error,omitempty"`
}

func (dsr *deleteStudentsResponse) toJSON() []byte {
	res, err := json.Marshal(dsr)
	if err != nil {
		logger.Logger.Error("failed to marshal delete students response", zap.Error(err))
		return nil
	}
	return res
}

func deleteStudents(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	var req deleteStudentsRequest
	var resp deleteStudentsResponse

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = &errorResponse{Code: http.StatusBadRequest, Message: "invalid request"}
		return
	}

	for _, id := range req.UserIDs {
		userID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			resp.Status = append(resp.Status, deleteStudentStatus{
				UserID:  id,
				Code:    http.StatusBadRequest,
				Message: "invalid user id",
			})
			continue
		}

		err = userService.DeleteByUserID(userID)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrUserNotFound):
				resp.Status = append(resp.Status, deleteStudentStatus{
					UserID:  id,
					Code:    http.StatusNotFound,
					Message: "user not found",
				})
			default:
				logger.Logger.Error("failed to delete user", zap.String("requestID", requestID), zap.String("userID", id), zap.Error(err))
				resp.Status = append(resp.Status, deleteStudentStatus{
					UserID:  id,
					Code:    http.StatusInternalServerError,
					Message: "failed to delete user",
				})
			}
			continue
		}

		resp.Status = append(resp.Status, deleteStudentStatus{
			UserID:  id,
			Code:    http.StatusOK,
			Message: "user deleted",
		})
		continue
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp.toJSON())
}
