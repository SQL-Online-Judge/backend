package restapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/SQL-Online-Judge/backend/internal/model"
	"github.com/SQL-Online-Judge/backend/internal/pkg/logger"
	"go.uber.org/zap"
)

type createProblemRequest struct {
	Title       string   `json:"title"`
	Tags        []string `json:"tags"`
	Content     string   `json:"content"`
	TimeLimit   int32    `json:"timeLimit"`
	MemoryLimit int32    `json:"memoryLimit"`
}

type createProblemResponse struct {
	ProblemID string         `json:"problemID,omitempty"`
	Error     *errorResponse `json:"error,omitempty"`
}

func (cpr *createProblemResponse) toJSON() []byte {
	res, err := json.Marshal(cpr)
	if err != nil {
		logger.Logger.Error("failed to marshal create class response", zap.Error(err))
		return nil
	}
	return res
}

func createProblem(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	var req createProblemRequest
	var resp createProblemResponse

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = &errorResponse{Code: http.StatusBadRequest, Message: "invalid request"}
		w.Write(resp.toJSON())
		return
	}

	authorID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		logger.Logger.Error("failed to get author id from context", zap.String("requestID", requestID), zap.Any("authorID", authorID))
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = &errorResponse{Code: http.StatusInternalServerError, Message: "failed to get author id from context"}
		w.Write(resp.toJSON())
		return
	}

	problem := model.NewProblem(&model.Problem{
		AuthorID:    authorID,
		Title:       req.Title,
		Tags:        req.Tags,
		Content:     req.Content,
		TimeLimit:   req.TimeLimit,
		MemoryLimit: req.MemoryLimit,
	})

	if !problem.IsValidProblem() {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = &errorResponse{Code: http.StatusBadRequest, Message: "invalid problem"}
		w.Write(resp.toJSON())
		return
	}

	problemID, err := problemService.CreateProblem(problem)
	if err != nil {
		logger.Logger.Error("failed to create problem", zap.String("requestID", requestID), zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = &errorResponse{Code: http.StatusInternalServerError, Message: "failed to create problem"}
		w.Write(resp.toJSON())
		return
	}

	resp.ProblemID = strconv.FormatInt(problemID, 10)
	w.Write(resp.toJSON())
}
