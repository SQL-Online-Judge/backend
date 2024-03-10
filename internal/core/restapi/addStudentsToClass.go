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

type addStudentsToClassRequest struct {
	Students []string `json:"students"`
}

type addStudentsToClassStatus struct {
	UserID  string `json:"userID"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type addStudentsToClassResponse struct {
	ClassID string                     `json:"classID,omitempty"`
	Status  []addStudentsToClassStatus `json:"status,omitempty"`
	Error   *errorResponse             `json:"error,omitempty"`
}

func (ascr *addStudentsToClassResponse) toJSON() []byte {
	res, err := json.Marshal(ascr)
	if err != nil {
		logger.Logger.Error("failed to marshal add students to class response", zap.Error(err))
		return nil
	}
	return res
}

func addStudentsToClass(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	var resp addStudentsToClassResponse

	sClassID := chi.URLParam(r, "classID")
	classID, err := strconv.ParseInt(sClassID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = &errorResponse{Code: http.StatusBadRequest, Message: "invalid class id"}
		w.Write(resp.toJSON())
		return
	}

	var req addStudentsToClassRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = &errorResponse{Code: http.StatusBadRequest, Message: "failed to decode request body"}
		w.Write(resp.toJSON())
		return
	}

	students := make([]int64, 0, len(req.Students))
	for _, student := range req.Students {
		studentID, err := strconv.ParseInt(student, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			resp.Error = &errorResponse{Code: http.StatusBadRequest, Message: "invalid student id"}
			w.Write(resp.toJSON())
			return
		}
		students = append(students, studentID)
	}

	teacherID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		logger.Logger.Error("failed to get teacher id from context", zap.String("requestID", requestID), zap.Any("teacherID", teacherID))
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = &errorResponse{Code: http.StatusInternalServerError, Message: "failed to get teacher id from context"}
		w.Write(resp.toJSON())
		return
	}

	status, err := classService.AddStudentsToClass(userService, teacherID, classID, students)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrClassNotFound):
			w.WriteHeader(http.StatusNotFound)
			resp.Error = &errorResponse{Code: http.StatusNotFound, Message: "class not found"}
		case errors.Is(err, service.ErrNotOfClassOwner):
			w.WriteHeader(http.StatusForbidden)
			resp.Error = &errorResponse{Code: http.StatusForbidden, Message: "not class owner"}
		default:
			logger.Logger.Error("failed to add students to class", zap.String("requestID", requestID), zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			resp.Error = &errorResponse{Code: http.StatusInternalServerError, Message: "failed to add students to class"}
		}
		w.Write(resp.toJSON())
		return
	}

	resp.ClassID = sClassID
	handleAddStudentsToClassStatus(status, &resp)

	w.WriteHeader(http.StatusOK)
	w.Write(resp.toJSON())
}

func handleAddStudentsToClassStatus(status map[int64]error, resp *addStudentsToClassResponse) {
	for studentID, err := range status {
		if err != nil {
			switch {
			case errors.Is(err, service.ErrUserNotFound):
				resp.Status = append(resp.Status, addStudentsToClassStatus{
					UserID: strconv.FormatInt(studentID, 10), Code: http.StatusNotFound, Message: "student not found"})
			case errors.Is(err, service.ErrUserNotStudent):
				resp.Status = append(resp.Status, addStudentsToClassStatus{
					UserID: strconv.FormatInt(studentID, 10), Code: http.StatusForbidden, Message: "user is not student"})
			case errors.Is(err, service.ErrStudentAlreadyInClass):
				resp.Status = append(resp.Status, addStudentsToClassStatus{
					UserID: strconv.FormatInt(studentID, 10), Code: http.StatusConflict, Message: "student already in class"})
			default:
				logger.Logger.Error("failed to add student to class", zap.Int64("studentID", studentID), zap.Error(err))
				resp.Status = append(resp.Status, addStudentsToClassStatus{
					UserID: strconv.FormatInt(studentID, 10), Code: http.StatusInternalServerError, Message: "failed to add student to class"})
			}
		} else {
			resp.Status = append(resp.Status, addStudentsToClassStatus{
				UserID: strconv.FormatInt(studentID, 10), Code: http.StatusOK, Message: "success"})
		}
	}
}
