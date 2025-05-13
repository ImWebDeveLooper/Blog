package posts

import (
	"errors"
	"time"
)

var (
	ErrPostNotFound = errors.New("post not found")
)

type Status string

const (
	StatusPublished Status = "Published"
	StatusPending   Status = "Pending"
	StatusDraft     Status = "Draft"
	StatusArchived  Status = "Archived"
)

func (s Status) String() string {
	return string(s)
}

type Post struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Author      string    `json:"author"`
	Slug        string    `json:"slug"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	PublishedAt time.Time `json:"published_at"`
}
