package model

import "time"

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
