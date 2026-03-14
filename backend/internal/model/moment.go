package model

import "time"

type Moment struct {
	ID        uint      `json:"id" db:"id"`
	Content   string    `json:"content" db:"content"`
	Images    []string  `json:"images" db:"images"` // 图片 URL 列表
	Video     string    `json:"video" db:"video"`   // 视频 URL（可选）
	LikeCount int       `json:"like_count" db:"like_count"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
