package model

import "time"

type Like struct {
	ID         uint      `json:"id" db:"id"`
	UserID     uint      `json:"user_id" db:"user_id"`
	TargetType string    `json:"target_type" db:"target_type"` // post | moment
	TargetID   uint      `json:"target_id" db:"target_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}
