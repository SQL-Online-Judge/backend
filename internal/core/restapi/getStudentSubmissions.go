package restapi

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/SQL-Online-Judge/backend/internal/pkg/logger"
	"go.uber.org/zap"
)

type submissionSummary struct {
	SubmissionID string `json:"submissionID"`
	SubmitTime   string `json:"submitTime"`
	TaskID       string `json:"taskID"`
	TaskName     string `json:"taskName"`
	ProblemID    string `json:"problemID"`
	ProblemTitle string `json:"problemTitle"`
	DBName       string `json:"dbName"`
	JudgeStatus  string `json:"judgeStatus"`
	TimeCost     int32  `json:"timeCost"`
}

type getStudentSubmissionsResponse struct {
	Submissions []*submissionSummary `json:"submissions,omitempty"`
	Error       *errorResponse       `json:"error,omitempty"`
}

func (gssr *getStudentSubmissionsResponse) toJSON() []byte {
	res, err := json.Marshal(gssr)
	if err != nil {
		logger.Logger.Error("failed to marshal get student submissions response", zap.Error(err))
		return nil
	}
	return res
}

func getStudentSubmissions(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	var resp getStudentSubmissionsResponse

	studentID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		logger.Logger.Error("failed to get user id from context", zap.String("requestID", requestID))
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = &errorResponse{Code: http.StatusInternalServerError, Message: "failed to get user id from context"}
		w.Write(resp.toJSON())
		return
	}

	submissions, err := submissionService.GetStudentSubmissions(studentID)
	if err != nil {
		logger.Logger.Error("failed to get student submissions", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = &errorResponse{Code: http.StatusInternalServerError, Message: "failed to get student submissions"}
		w.Write(resp.toJSON())
		return
	}

	resp.Submissions = make([]*submissionSummary, 0, len(submissions))
	for _, submission := range submissions {
		resp.Submissions = append(resp.Submissions, &submissionSummary{
			SubmissionID: strconv.FormatInt(submission.SubmissionID, 10),
			SubmitTime:   submission.SubmitTime.Format(time.RFC3339),
			TaskID:       strconv.FormatInt(submission.TaskID, 10),
			TaskName:     submission.TaskName,
			ProblemID:    strconv.FormatInt(submission.ProblemID, 10),
			ProblemTitle: submission.ProblemTitle,
			DBName:       submission.DBName,
			JudgeStatus:  submission.JudgeStatus,
			TimeCost:     submission.TimeCost,
		})
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp.toJSON())
}
