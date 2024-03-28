package model

import (
	"encoding/json"
	"fmt"
)

const (
	JudgeStatusPending           = "Pending"
	JudgeStatusQueued            = "Queued"
	JudgeStatusJudging           = "Judging"
	JudgeStatusAccepted          = "Accepted"
	JudgeStatusWrongAnswer       = "Wrong Answer"
	JudgeStatusTimeLimitExceeded = "Time Limit Exceeded"
	JudgeStatusRuntimeError      = "Runtime Error"
	JudgeStatusSystemError       = "System Error"
)

type JudgeSubmission struct {
	SubmissionID int64  `bson:"submissionID" json:"submissionID"`
	SubmittedSQL string `bson:"submittedSQL" json:"submittedSQL"`
}

type JudgeProblem struct {
	TimeLimit   int32 `bson:"timeLimit" json:"timeLimit"`
	MemoryLimit int32 `bson:"memoryLimit" json:"memoryLimit"`
}

type JudgeAnswer struct {
	DBName       string `bson:"dbName" json:"dbName"`
	PrepareSQL   string `bson:"prepareSQL" json:"prepareSQL"`
	AnswerSQL    string `bson:"answerSQL" json:"answerSQL"`
	JudgeSQL     string `bson:"judgeSQL" json:"judgeSQL"`
	AnswerOutput string `bson:"answerOutput" json:"answerOutput"`
}

type JudgeResult struct {
	JudgeStatus  string `bson:"judgeStatus" json:"judgeStatus"`
	TimeCost     int32  `bson:"timeCost" json:"timeCost"`
	JudgerOutput string `bson:"judgerOutput" json:"judgerOutput"`
}

type JudgeRequest struct {
	Submission *JudgeSubmission `json:"submission"`
	Problem    *JudgeProblem    `json:"problem"`
	Answer     *JudgeAnswer     `json:"answer"`
}

type JudgeResponse struct {
	SubmissionID int64        `json:"submissionID"`
	Result       *JudgeResult `json:"result"`
}

func (jr *JudgeRequest) ToJSON() (string, error) {
	j, err := json.Marshal(jr)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JudgeRequest: %w", err)
	}
	return string(j), nil
}

func (jr *JudgeRequest) FromJSON(j string) error {
	err := json.Unmarshal([]byte(j), jr)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JudgeRequest: %w", err)
	}
	return nil
}

func (jr *JudgeResponse) ToJSON() (string, error) {
	j, err := json.Marshal(jr)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JudgeResponse: %w", err)
	}
	return string(j), nil
}

func (jr *JudgeResponse) FromJSON(j string) error {
	err := json.Unmarshal([]byte(j), jr)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JudgeResponse: %w", err)
	}
	return nil
}
