package model

import (
	"time"
	"unicode/utf8"

	"github.com/SQL-Online-Judge/backend/internal/pkg/id"
)

type Submission struct {
	SubmissionID int64     `bson:"submissionID"`
	SubmitterID  int64     `bson:"submitterID"`
	SubmitTime   time.Time `bson:"submitTime"`
	TaskID       int64     `bson:"taskID"`
	ProblemID    int64     `bson:"problemID"`
	DBName       string    `bson:"dbName"`
	SubmittedSQL string    `bson:"submittedSQL"`
	JudgeStatus  string    `bson:"judgeStatus"`
	TimeCost     int32     `bson:"timeCost"`
	JudgerOutput string    `bson:"judgerOutput"`
}

func (s *Submission) IsValidDBName() bool {
	switch s.DBName {
	case "mysql", "opengauss":
		return true
	default:
		return false
	}
}

func (s *Submission) IsValidSubmittedSQL() bool {
	sqlLen := utf8.RuneCountInString(s.SubmittedSQL)
	return sqlLen >= 2 && sqlLen <= 65536
}

func (s *Submission) IsValidJudgeStatus() bool {
	switch s.JudgeStatus {
	case "Pending", "Queued", "Judging", "Accepted", "Wrong Answer", "Time Limit Exceeded", "Runtime Error", "System Error":
		return true
	default:
		return false
	}
}

func (s *Submission) IsValidSubmission() bool {
	return s.IsValidDBName() && s.IsValidSubmittedSQL() && s.IsValidJudgeStatus()
}

func NewSubmission(s *Submission) *Submission {
	return &Submission{
		SubmissionID: id.NewID(),
		SubmitterID:  s.SubmitterID,
		SubmitTime:   time.Now(),
		TaskID:       s.TaskID,
		ProblemID:    s.ProblemID,
		DBName:       s.DBName,
		SubmittedSQL: s.SubmittedSQL,
		JudgeStatus:  "Pending",
		TimeCost:     0,
		JudgerOutput: "",
	}
}

type SubmissionSummary struct {
	SubmissionID int64     `bson:"submissionID"`
	SubmitTime   time.Time `bson:"submitTime"`
	TaskID       int64     `bson:"taskID"`
	TaskName     string    `bson:"taskName"`
	ProblemID    int64     `bson:"problemID"`
	ProblemTitle string    `bson:"problemTitle"`
	DBName       string    `bson:"dbName"`
	JudgeStatus  string    `bson:"judgeStatus"`
	TimeCost     int32     `bson:"timeCost"`
}

type SubmitedSQL struct {
	SubmissionID int64  `bson:"submissionID"`
	SubmittedSQL string `bson:"submittedSQL"`
}
