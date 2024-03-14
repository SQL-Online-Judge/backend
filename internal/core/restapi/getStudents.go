package restapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/SQL-Online-Judge/backend/internal/model"
	"github.com/SQL-Online-Judge/backend/internal/pkg/logger"
	"go.uber.org/zap"
)

type student struct {
	UserID   string `json:"userID"`
	Username string `json:"username"`
}

func newStudentFromModel(s *model.User) *student {
	return &student{
		UserID:   strconv.FormatInt(s.UserID, 10),
		Username: s.Username,
	}
}

type getStudentsResponse struct {
	Students []*student     `json:"students,omitempty"`
	Error    *errorResponse `json:"error,omitempty"`
}

func (gsr *getStudentsResponse) toJSON() []byte {
	res, err := json.Marshal(gsr)
	if err != nil {
		logger.Logger.Error("failed to marshal get students response", zap.Error(err))
		return nil
	}
	return res
}

func getStudents(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	var resp getStudentsResponse

	params := r.URL.Query()
	contains := params.Get("contains")
	students, err := userService.GetStudents(contains)
	if err != nil {
		logger.Logger.Error("failed to get students",
			zap.String("requestID", requestID),
			zap.String("contains", contains),
			zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = &errorResponse{Code: http.StatusInternalServerError, Message: "failed to get students"}
		w.Write(resp.toJSON())
		return
	}

	if len(students) == 0 {
		w.WriteHeader(http.StatusNotFound)
		resp.Error = &errorResponse{Code: http.StatusNotFound, Message: "students not found"}
		w.Write(resp.toJSON())
		return
	}

	resp.Students = make([]*student, 0, len(students))
	for _, s := range students {
		resp.Students = append(resp.Students, newStudentFromModel(s))
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp.toJSON())
}
