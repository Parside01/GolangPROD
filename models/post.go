package models

import "time"

type Post struct {
	ID            string    `json:"id" mapstructure:"id" db:"id"`
	Author        string    `json:"author" mapstructure:"author" db:"author"`
	Content       string    `json:"content" mapstructure:"content" db:"content"`
	Tags          []string  `json:"tags" mapstructure:"tags" db:"tags"`
	LikesCount    int       `json:"likesCount" mapstructure:"likesCount" db:"likesCount"`
	DislikesCount int       `json:"dislikesCount" mapstructure:"dislikesCount" db:"dislikesCount"`
	CreatedAt     time.Time `json:"createdAt" mapstructure:"createdAt" db:"createdAt"`
}

type MainContentPost struct {
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

func (p *Post) GetMainContent() *MainContentPost {
	return &MainContentPost{
		Content: p.Content,
		Tags:    p.Tags,
	}
}
