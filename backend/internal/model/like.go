package model

type Like struct {
	ID         uint   `gorm:"primaryKey"`
	UserID     uint   `gorm:"uniqueIndex:idx_like"`
	TargetType string `gorm:"uniqueIndex:idx_like"`
	TargetID   uint   `gorm:"uniqueIndex:idx_like"`
}
