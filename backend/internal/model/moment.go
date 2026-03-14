package model

import (
	"time"

	"github.com/lib/pq"
)

type Moment struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Content   string         `json:"content"`
	Images    pq.StringArray `gorm:"type:text[]" json:"images"`
	Video     string         `json:"video"`
	LikeCount int            `gorm:"default:0" json:"like_count"`
	CreatedAt time.Time      `json:"created_at"`
}
