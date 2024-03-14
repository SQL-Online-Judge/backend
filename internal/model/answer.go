package model

import (
	"unicode/utf8"

	"github.com/SQL-Online-Judge/backend/internal/pkg/id"
)

type Answer struct {
	AnswerID     int64  `bson:"answerID"`
	ProblemID    int64  `bson:"problemID"`
	DBName       string `bson:"dbName"`
	PrepareSQL   string `bson:"prepareSQL"`
	AnswerSQL    string `bson:"answerSQL"`
	JudgeSQL     string `bson:"judgeSQL"`
	AnswerOutput string `bson:"answerOutput"`
	IsReady      bool   `bson:"isReady"`
	ImageName    string `bson:"imageName"`
	Deleted      bool   `bson:"deleted"`
}

func (a *Answer) IsValidDBName() bool {
	switch a.DBName {
	case "mysql", "opengauss":
		return true
	default:
		return false
	}
}

func (a *Answer) IsValidPrepareSQL() bool {
	sqlLen := utf8.RuneCountInString(a.PrepareSQL)
	return sqlLen >= 2 && sqlLen <= 65536
}

func (a *Answer) IsValidAnswerSQL() bool {
	sqlLen := utf8.RuneCountInString(a.AnswerSQL)
	return sqlLen >= 2 && sqlLen <= 65536
}

func (a *Answer) IsValidAnswer() bool {
	return a.IsValidDBName() && a.IsValidPrepareSQL() && a.IsValidAnswerSQL()
}

func NewAnswer(a *Answer) *Answer {
	return &Answer{
		AnswerID:     id.NewID(),
		ProblemID:    a.ProblemID,
		DBName:       a.DBName,
		PrepareSQL:   a.PrepareSQL,
		AnswerSQL:    a.AnswerSQL,
		JudgeSQL:     a.JudgeSQL,
		AnswerOutput: a.AnswerOutput,
		IsReady:      a.IsReady,
		ImageName:    a.ImageName,
		Deleted:      a.Deleted,
	}
}
