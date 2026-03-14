package repo

import (
	"fmt"

	"github.com/isyayanoaa/go-blog/internal/model"
	"github.com/isyayanoaa/go-blog/pkg/database"
)

func AddLike(userID uint, targetType string, targetID int) error {
	like := model.Like{UserID: userID, TargetType: targetType, TargetID: uint(targetID)}
	database.DB.Where(like).FirstOrCreate(&like)
	return database.DB.Exec(fmt.Sprintf("UPDATE %s SET like_count = like_count + 1 WHERE id = ?", likeTable(targetType)), targetID).Error
}

func RemoveLike(userID uint, targetType string, targetID int) error {
	database.DB.Where("user_id = ? AND target_type = ? AND target_id = ?", userID, targetType, targetID).Delete(&model.Like{})
	return database.DB.Exec(fmt.Sprintf("UPDATE %s SET like_count = GREATEST(like_count - 1, 0) WHERE id = ?", likeTable(targetType)), targetID).Error
}

func likeTable(targetType string) string {
	if targetType == "moment" { return "moments" }
	return "posts"
}
