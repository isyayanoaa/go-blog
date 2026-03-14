package model

import "time"

type Comment struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `json:"user_id"`
	TargetType string    `json:"target_type"`
	TargetID   uint      `json:"target_id"`
	ParentID   *uint     `json:"parent_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}
