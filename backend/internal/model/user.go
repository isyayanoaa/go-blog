package model

import "time"

type User struct {
	ID        uint      `json:"id" db:"id"`
	GithubID  int64     `json:"github_id" db:"github_id"`
	Username  string    `json:"username" db:"username"`
	Avatar    string    `json:"avatar" db:"avatar"`
	IsAdmin   bool      `json:"is_admin" db:"is_admin"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
