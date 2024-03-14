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

type answer struct {
	AnswerID     string `json:"answerID"`
	DBName       string `json:"dbName"`
	PrepareSQL   string `json:"prepareSQL"`
	AnswerSQL    string `json:"answerSQL"`
	JudgeSQL     string `json:"judgeSQL"`
	AnswerOutput string `json:"answerOutput"`
	IsReady      bool   `json:"isReady"`
}

type getAnswersResponse struct {
	ProblemID string         `json:"problemID,omitempty"`
	Answers   []*answer      `json:"answers,omitempty"`
	Error     *errorResponse `json:"error,omitempty"`
}

func (gar *getAnswersResponse) toJSON() []byte {
	res, err := json.Marshal(gar)
	if err != nil {
		logger.Logger.Error("failed to marshal get answers response", zap.Error(err))
		return nil
	}

	return res
}

func (gar *getAnswersResponse) fromAnswers(answers []*model.Answer) {
	gar.Answers = make([]*answer, 0, len(answers))
	for _, a := range answers {
		gar.Answers = append(gar.Answers, &answer{
			AnswerID:     strconv.FormatInt(a.AnswerID, 10),
			DBName:       a.DBName,
			PrepareSQL:   a.PrepareSQL,
			AnswerSQL:    a.AnswerSQL,
			JudgeSQL:     a.JudgeSQL,
			AnswerOutput: a.AnswerOutput,
			IsReady:      a.IsReady,
		})
	}
}

func getAnswers(w http.ResponseWriter, r *http.Request) {
	var resp getAnswersResponse

	problemID, err := strconv.ParseInt(chi.URLParam(r, "problemID"), 10, 64)
	if err != nil {
		resp.Error = &errorResponse{Code: http.StatusBadRequest, Message: "invalid problem id"}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp.toJSON())
		return
	}

	answers, err := answerService.GetAnswers(problemService, problemID)
	if err != nil {
		handleGetAnswersError(w, &resp, err)
		return
	}

	if len(answers) == 0 {
		resp.Error = &errorResponse{Code: http.StatusNotFound, Message: "no answers found"}
		w.WriteHeader(http.StatusNotFound)
		w.Write(resp.toJSON())
		return
	}

	resp.ProblemID = strconv.FormatInt(problemID, 10)
	resp.fromAnswers(answers)
	w.WriteHeader(http.StatusOK)
	w.Write(resp.toJSON())
}

func handleGetAnswersError(w http.ResponseWriter, resp *getAnswersResponse, err error) {
	switch {
	case errors.Is(err, service.ErrProblemNotFound):
		resp.Error = &errorResponse{Code: http.StatusNotFound, Message: "problem not found"}
		w.WriteHeader(http.StatusNotFound)
	default:
		logger.Logger.Error("failed to get answers", zap.Error(err))
		resp.Error = &errorResponse{Code: http.StatusInternalServerError, Message: "failed to get answers"}
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write(resp.toJSON())
}
