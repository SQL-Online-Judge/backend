package restapi

import (
	"net/http"

	"github.com/SQL-Online-Judge/backend/internal/pkg/logger"
	"go.uber.org/zap"
)

func getTeacherProblems(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	var resp getProblemsResponse

	teacherID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		logger.Logger.Error("failed to get teacher id from context", zap.String("requestID", requestID), zap.Any("authorID", teacherID))
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = &errorResponse{Code: http.StatusInternalServerError, Message: "failed to get author id from context"}
		w.Write(resp.toJSON())
		return
	}

	problems, err := problemService.GetTeacherProblems(teacherID)
	if err != nil {
		logger.Logger.Error("failed to get problems",
			zap.String("requestID", requestID),
			zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = &errorResponse{Code: http.StatusInternalServerError, Message: "failed to get problems"}
		w.Write(resp.toJSON())
		return
	}

	if len(problems) == 0 {
		w.WriteHeader(http.StatusNotFound)
		resp.Error = &errorResponse{Code: http.StatusNotFound, Message: "problems not found"}
		w.Write(resp.toJSON())
		return
	}

	resp.Problems = make([]*problem, 0, len(problems))
	for _, p := range problems {
		resp.Problems = append(resp.Problems, newProblemFromModel(p))
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp.toJSON())
}
