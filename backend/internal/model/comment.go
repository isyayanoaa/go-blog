package model

import "time"

type Comment struct {
	ID         uint      `json:"id" db:"id"`
	UserID     uint      `json:"user_id" db:"user_id"`
	User       *User     `json:"user,omitempty"`
	TargetType string    `json:"target_type" db:"target_type"` // post | moment
	TargetID   uint      `json:"target_id" db:"target_id"`
	ParentID   *uint     `json:"parent_id" db:"parent_id"` // 一级回复，nil 表示顶层评论
	Content    string    `json:"content" db:"content"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}
