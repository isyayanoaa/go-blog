package dto

import "time"

type CreateCommentReq struct {
	TargetType string `json:"target_type" binding:"required"`
	TargetID   uint   `json:"target_id"   binding:"required"`
	ParentID   *uint  `json:"parent_id"`
	Content    string `json:"content"     binding:"required"`
}

type CommentResp struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id"`
	Username   string    `json:"username"`
	Avatar     string    `json:"avatar"`
	TargetType string    `json:"target_type"`
	TargetID   uint      `json:"target_id"`
	ParentID   *uint     `json:"parent_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}
