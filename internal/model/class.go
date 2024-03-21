package model

import (
	"unicode/utf8"

	"github.com/SQL-Online-Judge/backend/internal/pkg/id"
)

type Class struct {
	ClassID   int64   `bson:"classID"`
	ClassName string  `bson:"className"`
	TeacherID int64   `bson:"teacherID"`
	Students  []int64 `bson:"students"`
	Tasks     []int64 `bson:"tasks"`
	Deleted   bool    `bson:"deleted"`
}

func NewClass(className string, teacherID int64) *Class {
	return &Class{
		ClassID:   id.NewID(),
		ClassName: className,
		TeacherID: teacherID,
		Students:  []int64{},
		Tasks:     []int64{},
	}
}

func (c *Class) IsValidClassName() bool {
	nameLen := utf8.RuneCountInString(c.ClassName)
	return nameLen >= 2 && nameLen <= 32
}
