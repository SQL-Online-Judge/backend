package model

import (
	"unicode/utf8"

	"github.com/SQL-Online-Judge/backend/internal/pkg/id"
)

type Problem struct {
	ProblemID   int64    `bson:"problemID"`
	AuthorID    int64    `bson:"authorID"`
	Title       string   `bson:"title"`
	Tags        []string `bson:"tags"`
	Content     string   `bson:"content"`
	TimeLimit   int32    `bson:"timeLimit"`
	MemoryLimit int32    `bson:"memoryLimit"`
	Deleted     bool     `bson:"deleted"`
}

func (p *Problem) IsValidTitle() bool {
	titleLen := utf8.RuneCountInString(p.Title)
	return titleLen >= 2 && titleLen <= 64
}

func (p *Problem) IsValidTags() bool {
	if p.Tags == nil {
		return false
	}
	for _, tag := range p.Tags {
		tagLen := utf8.RuneCountInString(tag)
		if tagLen < 2 || tagLen > 32 {
			return false
		}
	}
	return true
}

func (p *Problem) IsValidContent() bool {
	contentLen := utf8.RuneCountInString(p.Content)
	return contentLen >= 2 && contentLen <= 65536
}

func (p *Problem) IsValidTimeLimit() bool {
	return p.TimeLimit >= 100 && p.TimeLimit <= 60000
}

func (p *Problem) IsValidMemoryLimit() bool {
	return p.MemoryLimit >= 200 && p.MemoryLimit <= 4096
}

func (p *Problem) IsValidProblem() bool {
	return p.IsValidTitle() && p.IsValidTags() && p.IsValidContent() && p.IsValidTimeLimit() && p.IsValidMemoryLimit()
}

func NewProblem(p *Problem) *Problem {
	return &Problem{
		ProblemID:   id.NewID(),
		AuthorID:    p.AuthorID,
		Title:       p.Title,
		Tags:        p.Tags,
		Content:     p.Content,
		TimeLimit:   p.TimeLimit,
		MemoryLimit: p.MemoryLimit,
		Deleted:     false,
	}
}
