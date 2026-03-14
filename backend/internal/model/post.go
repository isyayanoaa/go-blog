package model

import "time"

type Post struct {
	ID          uint      `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Slug        string    `json:"slug" db:"slug"`
	Content     string    `json:"content" db:"content"`
	Summary     string    `json:"summary" db:"summary"`
	Cover       string    `json:"cover" db:"cover"`
	Category    string    `json:"category" db:"category"`
	Tags        []string  `json:"tags" db:"tags"`
	Status      string    `json:"status" db:"status"` // published | draft
	LikeCount   int       `json:"like_count" db:"like_count"`
	ViewCount   int       `json:"view_count" db:"view_count"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
