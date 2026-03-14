package repo

import (
	"github.com/isyayanoaa/go-blog/internal/dto"
	"github.com/isyayanoaa/go-blog/internal/model"
	"github.com/isyayanoaa/go-blog/pkg/database"
)

func ListComments(targetType, targetID string) ([]dto.CommentResp, error) {
	var result []dto.CommentResp
	err := database.DB.Model(&model.Comment{}).
		Select("comments.id, comments.user_id, users.username, users.avatar, comments.target_type, comments.target_id, comments.parent_id, comments.content, comments.created_at").
		Joins("JOIN users ON users.id = comments.user_id").
		Where("comments.target_type = ? AND comments.target_id = ?", targetType, targetID).
		Order("comments.created_at ASC").Scan(&result).Error
	return result, err
}

func CreateComment(userID uint, req dto.CreateCommentReq) (uint, error) {
	c := model.Comment{
		UserID: userID, TargetType: req.TargetType,
		TargetID: req.TargetID, ParentID: req.ParentID, Content: req.Content,
	}
	err := database.DB.Create(&c).Error
	return c.ID, err
}

func DeleteComment(id string) error {
	return database.DB.Delete(&model.Comment{}, id).Error
}
