package dto

type LikeReq struct {
	TargetType string `json:"target_type" binding:"required"`
	TargetID   int    `json:"target_id"   binding:"required"`
}
