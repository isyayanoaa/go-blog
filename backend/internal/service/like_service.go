package service

import (
	"github.com/isyayanoaa/go-blog/internal/dto"
	"github.com/isyayanoaa/go-blog/internal/repo"
)

func AddLike(userID uint, req dto.LikeReq) error {
	return repo.AddLike(userID, req.TargetType, req.TargetID)
}
func RemoveLike(userID uint, req dto.LikeReq) error {
	return repo.RemoveLike(userID, req.TargetType, req.TargetID)
}
