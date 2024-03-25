package model

import "time"

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
