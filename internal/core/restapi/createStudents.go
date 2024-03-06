package restapi

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/SQL-Online-Judge/backend/internal/core/service"
	"github.com/SQL-Online-Judge/backend/internal/model"
	"github.com/SQL-Online-Judge/backend/internal/pkg/logger"
	"go.uber.org/zap"
)

type createStudentsRequest struct {
	Usernames []string `json:"usernames"`
}

func (csr *createStudentsRequest) toUsers() []*model.User {
	users := make([]*model.User, 0, len(csr.Usernames))
	for _, username := range csr.Usernames {
		h := sha256.New()
		h.Write([]byte(username))
		p := h.Sum(nil)
		users = append(users, &model.User{
			Username: username,
			Password: hex.EncodeToString(p),
			Role:     "student",
		})
	}
	return users
}

type createStudentStatus struct {
	UserID   string `json:"userID,omitempty"`
	Username string `json:"username"`
	Code     int    `json:"code"`
	Message  string `json:"message"`
}

type createStudentsResponse struct {
	Status []createStudentStatus `json:"status,omitempty"`
	Error  *errorResponse        `json:"error,omitempty"`
}

func (csr *createStudentsResponse) toJSON() []byte {
	res, err := json.Marshal(csr)
	if err != nil {
		logger.Logger.Error("failed to marshal create students response", zap.Error(err))
		return nil
	}
	return res
}

func createStudents(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	var req createStudentsRequest
	var students []*model.User
	var resp createStudentsResponse

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = &errorResponse{Code: http.StatusBadRequest, Message: "invalid request"}
		w.Write(resp.toJSON())
		return
	}

	students = req.toUsers()
	for _, student := range students {
		if !student.IsValidUser() {
			resp.Status = append(resp.Status, createStudentStatus{
				Username: student.Username,
				Code:     http.StatusBadRequest,
				Message:  "invalid username or password",
			})
			continue
		}

		userID, err := userService.CreateUser(student.Username, student.Password, student.Role)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrUserConflict):
				resp.Status = append(resp.Status, createStudentStatus{
					Username: student.Username,
					Code:     http.StatusConflict,
					Message:  "username already exists",
				})
			default:
				logger.Logger.Error("failed to create user", zap.String("requestID", requestID), zap.String("username", student.Username), zap.Error(err))
				resp.Status = append(resp.Status, createStudentStatus{
					Username: student.Username,
					Code:     http.StatusInternalServerError,
					Message:  "failed to create user",
				})
			}
			continue
		}

		resp.Status = append(resp.Status, createStudentStatus{
			UserID:   strconv.FormatInt(userID, 10),
			Username: student.Username,
			Code:     http.StatusCreated,
			Message:  "user created",
		})
		continue
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp.toJSON())
}
