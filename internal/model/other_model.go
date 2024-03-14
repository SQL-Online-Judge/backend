package model

import "time"

type ProblemSet struct {
	ProblemSetID   int64   `bson:"problemSetID"`
	AuthorID       int64   `bson:"authorID"`
	ProblemSetName string  `bson:"problemSetName"`
	IsPublic       bool    `bson:"isPublic"`
	Problems       []int64 `bson:"problems"`
}

type ClassProblemSet struct {
	ClassID     int64   `bson:"classID"`
	ProblemSets []int64 `bson:"problemSets"`
}

type Task struct {
	TaskID        int64     `bson:"taskID"`
	AuthorID      int64     `bson:"authorID"`
	Classes       []int64   `bson:"classes"`
	Problems      []int64   `bson:"problems"`
	IsTimeLimited bool      `bson:"isTimeLimited"`
	BeginTime     time.Time `bson:"beginTime"`
	EndTime       time.Time `bson:"endTime"`
}

type Submission struct {
	SubmissionID int64     `bson:"submissionID"`
	SubmitterID  int64     `bson:"submitterID"`
	SubmitTime   time.Time `bson:"submitTime"`
	ProblemID    int64     `bson:"problemID"`
	TaskID       int64     `bson:"taskID"`
	DBName       string    `bson:"dbName"`
	SubmittedSQL string    `bson:"submittedSQL"`
	JudgeStatus  string    `bson:"judgeStatus"`
	TimeCost     int32     `bson:"timeCost"`
	JudgerOutput string    `bson:"judgerOutput"`
}

type Message struct {
	MessageID int64     `bson:"messageID"`
	SenderID  int64     `bson:"senderID"`
	Title     string    `bson:"title"`
	Content   string    `bson:"content"`
	SendTime  time.Time `bson:"sendTime"`
}

type MessageRead struct {
	MessageID int64 `bson:"messageID"`
	IsRead    bool  `bson:"isRead"`
}

type MessageBox struct {
	UserID   int64       `bson:"userID"`
	Messages MessageRead `bson:"messages"`
}
