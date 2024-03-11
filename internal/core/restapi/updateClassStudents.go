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

type updateClassStudentsRequest struct {
	Students []string `json:"students"`
}

type updateClassStudentsStatus struct {
	UserID  string `json:"userID"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type updateClassStudentsResponse struct {
	ClassID string                      `json:"classID,omitempty"`
	Status  []updateClassStudentsStatus `json:"status,omitempty"`
	Error   *errorResponse              `json:"error,omitempty"`
}

func (ucsr *updateClassStudentsResponse) toJSON() []byte {
	res, err := json.Marshal(ucsr)
	if err != nil {
		logger.Logger.Error("failed to marshal update class students response", zap.Error(err))
		return nil
	}
	return res
}

func updateClassStudents(w http.ResponseWriter, r *http.Request, updateType string) {
	requestID := getRequestID(r)
	var resp updateClassStudentsResponse

	sClassID := chi.URLParam(r, "classID")
	classID, err := strconv.ParseInt(sClassID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = &errorResponse{Code: http.StatusBadRequest, Message: "invalid class id"}
		w.Write(resp.toJSON())
		return
	}

	var req updateClassStudentsRequest
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

	var status map[int64]error
	switch updateType {
	case "add":
		status, err = classService.AddStudentsToClass(userService, teacherID, classID, students)
	case "remove":
		status, err = classService.RemoveStudentsFromClass(userService, teacherID, classID, students)
	default:
		logger.Logger.Error("invalid update type", zap.String("requestID", requestID), zap.String("updateType", updateType))
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = &errorResponse{Code: http.StatusInternalServerError, Message: "invalid update type"}
		w.Write(resp.toJSON())
		return
	}

	if err != nil {
		resp.Error = handleUpdateClassStudentsError(w, err)
		w.Write(resp.toJSON())
		return
	}

	resp.ClassID = sClassID
	handleUpdateClassStudentsStatus(status, &resp)
	w.WriteHeader(http.StatusOK)
	w.Write(resp.toJSON())
}

func handleUpdateClassStudentsError(w http.ResponseWriter, err error) *errorResponse {
	var resp *errorResponse
	switch {
	case errors.Is(err, service.ErrClassNotFound):
		w.WriteHeader(http.StatusNotFound)
		resp = &errorResponse{Code: http.StatusNotFound, Message: "class not found"}
	case errors.Is(err, service.ErrNotOfClassOwner):
		w.WriteHeader(http.StatusForbidden)
		resp = &errorResponse{Code: http.StatusForbidden, Message: "not class owner"}
	default:
		logger.Logger.Error("failed to update class students", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		resp = &errorResponse{Code: http.StatusInternalServerError, Message: "failed to update class students"}
	}
	return resp
}

func handleUpdateClassStudentsStatus(status map[int64]error, resp *updateClassStudentsResponse) {
	for studentID, err := range status {
		if err == nil {
			resp.Status = append(resp.Status, updateClassStudentsStatus{
				UserID: strconv.FormatInt(studentID, 10), Code: http.StatusOK, Message: "success"})
			continue
		}
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			resp.Status = append(resp.Status, updateClassStudentsStatus{
				UserID: strconv.FormatInt(studentID, 10), Code: http.StatusNotFound, Message: "student not found"})
		case errors.Is(err, service.ErrUserNotStudent):
			resp.Status = append(resp.Status, updateClassStudentsStatus{
				UserID: strconv.FormatInt(studentID, 10), Code: http.StatusForbidden, Message: "user is not student"})
		case errors.Is(err, service.ErrStudentNotInClass):
			resp.Status = append(resp.Status, updateClassStudentsStatus{
				UserID: strconv.FormatInt(studentID, 10), Code: http.StatusNotFound, Message: "student not in class"})
		case errors.Is(err, service.ErrStudentAlreadyInClass):
			resp.Status = append(resp.Status, updateClassStudentsStatus{
				UserID: strconv.FormatInt(studentID, 10), Code: http.StatusConflict, Message: "student already in class"})
		default:
			logger.Logger.Error("failed to update class students", zap.Int64("studentID", studentID), zap.Error(err))
			resp.Status = append(resp.Status, updateClassStudentsStatus{
				UserID: strconv.FormatInt(studentID, 10), Code: http.StatusInternalServerError, Message: "failed to update class students"})
		}
	}
}
