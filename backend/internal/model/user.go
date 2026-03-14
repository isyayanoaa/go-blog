package model

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	GithubID int64  `gorm:"uniqueIndex" json:"github_id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	IsAdmin  bool   `gorm:"default:false" json:"is_admin"`
}
