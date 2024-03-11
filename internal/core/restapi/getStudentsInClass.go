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

type getStudentInClassResponse struct {
	ClassID  string         `json:"classID,omitempty"`
	Students []*student     `json:"students"`
	Error    *errorResponse `json:"error,omitempty"`
}

func (gscr *getStudentInClassResponse) toJSON() []byte {
	res, err := json.Marshal(gscr)
	if err != nil {
		logger.Logger.Error("failed to marshal get students in class response", zap.Error(err))
		return nil
	}
	return res
}

func (gscr *getStudentInClassResponse) fromUsers(students []*model.User) {
	gscr.Students = make([]*student, 0, len(students))
	for _, stu := range students {
		gscr.Students = append(gscr.Students, &student{UserID: strconv.FormatInt(stu.UserID, 10), Username: stu.Username})
	}
}

func getStudentsInClass(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	var resp getStudentInClassResponse

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

	students, err := classService.GetStudentsInClass(teacherID, classID)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		resp.ClassID = sClassID
		resp.fromUsers(students)
		w.Write(resp.toJSON())
		return
	}

	switch {
	case errors.Is(err, service.ErrClassNotFound):
		w.WriteHeader(http.StatusNotFound)
		resp.Error = &errorResponse{Code: http.StatusNotFound, Message: "class not found"}
	case errors.Is(err, service.ErrNotOfClassOwner):
		w.WriteHeader(http.StatusForbidden)
		resp.Error = &errorResponse{Code: http.StatusForbidden, Message: "not class owner"}
	default:
		logger.Logger.Error("failed to get students in class", zap.String("requestID", requestID), zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = &errorResponse{Code: http.StatusInternalServerError, Message: "failed to get students in class"}
	}
	w.Write(resp.toJSON())
}
