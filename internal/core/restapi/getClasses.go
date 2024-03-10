package restapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/SQL-Online-Judge/backend/internal/pkg/logger"
	"go.uber.org/zap"
)

type classIDAndName struct {
	ClassID   string `json:"classID"`
	ClassName string `json:"className"`
}

type getClassesResponse struct {
	Classes []classIDAndName `json:"classes,omitempty"`
	Error   *errorResponse   `json:"error,omitempty"`
}

func (gcr *getClassesResponse) toJSON() []byte {
	res, err := json.Marshal(gcr)
	if err != nil {
		logger.Logger.Error("failed to marshal get classes response", zap.Error(err))
		return nil
	}
	return res
}

func getClasses(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	var resp getClassesResponse

	teacherID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		logger.Logger.Error("failed to get teacher id", zap.String("requestID", requestID))
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = &errorResponse{Code: http.StatusInternalServerError, Message: "failed to get teacher id"}
		w.Write(resp.toJSON())
		return
	}

	classes, err := classService.GetClasses(teacherID)
	if err != nil {
		logger.Logger.Error("failed to get classes", zap.String("requestID", requestID), zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = &errorResponse{Code: http.StatusInternalServerError, Message: "failed to get classes"}
		w.Write(resp.toJSON())
		return
	}

	if len(classes) == 0 {
		w.WriteHeader(http.StatusNotFound)
		resp.Error = &errorResponse{Code: http.StatusNotFound, Message: "no classes found"}
		w.Write(resp.toJSON())
		return
	}

	for _, class := range classes {
		resp.Classes = append(resp.Classes, classIDAndName{ClassID: strconv.FormatInt(class.ClassID, 10), ClassName: class.ClassName})
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp.toJSON())
}
