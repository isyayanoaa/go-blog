package service

import (
	"github.com/isyayanoaa/go-blog/internal/dto"
	"github.com/isyayanoaa/go-blog/internal/repo"
)

func ListComments(targetType, targetID string) ([]dto.CommentResp, error) {
	return repo.ListComments(targetType, targetID)
}
func CreateComment(userID uint, req dto.CreateCommentReq) (uint, error) {
	return repo.CreateComment(userID, req)
}
func DeleteComment(id string) error { return repo.DeleteComment(id) }
