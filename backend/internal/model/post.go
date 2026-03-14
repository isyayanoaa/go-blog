package model

import (
	"time"

	"github.com/lib/pq"
)

// Post 是数据库实体，对应 posts 表
type Post struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `gorm:"not null" json:"title"`
	Slug      string         `gorm:"uniqueIndex" json:"slug"`
	Content   string         `json:"content"`
	Summary   string         `json:"summary"`
	Cover     string         `json:"cover"`
	Category  string         `json:"category"`
	Tags      pq.StringArray `gorm:"type:text[]" json:"tags"`
	Status    string         `gorm:"default:draft" json:"status"`
	LikeCount int            `gorm:"default:0" json:"like_count"`
	ViewCount int            `gorm:"default:0" json:"view_count"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}
