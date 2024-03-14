package restapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/SQL-Online-Judge/backend/internal/model"
	"github.com/SQL-Online-Judge/backend/internal/pkg/logger"
	"go.uber.org/zap"
)

type problem struct {
	ProblemID string   `json:"problemID"`
	Title     string   `json:"title"`
	Tags      []string `json:"tags"`
}

func newProblemFromModel(p *model.Problem) *problem {
	return &problem{
		ProblemID: strconv.FormatInt(p.ProblemID, 10),
		Title:     p.Title,
		Tags:      p.Tags,
	}
}

type getProblemsResponse struct {
	Problems []*problem     `json:"problems,omitempty"`
	Error    *errorResponse `json:"error,omitempty"`
}

func (gpr *getProblemsResponse) toJSON() []byte {
	res, err := json.Marshal(gpr)
	if err != nil {
		logger.Logger.Error("failed to marshal get problems response", zap.Error(err))
		return nil
	}
	return res
}

func getProblems(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	var resp getProblemsResponse

	params := r.URL.Query()
	contains := params.Get("contains")

	problems, err := problemService.GetProblems(contains)
	if err != nil {
		logger.Logger.Error("failed to get problems",
			zap.String("requestID", requestID),
			zap.String("contains", contains),
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
