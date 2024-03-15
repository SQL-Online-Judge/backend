package model

import (
	"time"
	"unicode/utf8"

	"github.com/SQL-Online-Judge/backend/internal/pkg/id"
)

type TaskProblem struct {
	ProblemID int64   `bson:"problemID"`
	Score     float64 `bson:"score"`
}

func (tp *TaskProblem) IsValidScore() bool {
	return tp.Score > 0.0
}

type Task struct {
	TaskID        int64          `bson:"taskID"`
	AuthorID      int64          `bson:"authorID"`
	TaskName      string         `bson:"taskName"`
	Problems      []*TaskProblem `bson:"problems"`
	IsTimeLimited bool           `bson:"isTimeLimited"`
	BeginTime     time.Time      `bson:"beginTime"`
	EndTime       time.Time      `bson:"endTime"`
	Deleted       bool           `bson:"deleted"`
}

func (t *Task) IsValidTaskName() bool {
	nameLen := utf8.RuneCountInString(t.TaskName)
	return nameLen >= 2 && nameLen <= 64
}

func (t *Task) IsValidProblems() bool {
	if len(t.Problems) == 0 {
		return true
	}
	for _, problem := range t.Problems {
		if problem.Score < 0.0 {
			return false
		}
	}
	return true
}

func (t *Task) IsValidTime() bool {
	if !t.IsTimeLimited {
		t.BeginTime = time.Time{}
		t.EndTime = time.Time{}
		return true
	}
	return t.BeginTime.Before(t.EndTime)
}

func (t *Task) IsValidTask() bool {
	return t.IsValidTaskName() && t.IsValidProblems() && t.IsValidTime()
}

func NewTask(t *Task) *Task {
	return &Task{
		TaskID:        id.NewID(),
		AuthorID:      t.AuthorID,
		TaskName:      t.TaskName,
		Problems:      []*TaskProblem{},
		IsTimeLimited: t.IsTimeLimited,
		BeginTime:     t.BeginTime,
		EndTime:       t.EndTime,
		Deleted:       false,
	}
}
