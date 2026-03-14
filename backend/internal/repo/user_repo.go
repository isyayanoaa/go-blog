package repo

import (
	"github.com/isyayanoaa/go-blog/internal/model"
	"github.com/isyayanoaa/go-blog/pkg/database"
)

func UpsertUser(u model.User) (uint, error) {
	result := database.DB.Where(model.User{GithubID: u.GithubID}).FirstOrCreate(&u)
	if result.Error != nil {
		return 0, result.Error
	}
	// 更新最新的 username/avatar/is_admin
	database.DB.Model(&u).Updates(model.User{
		Username: u.Username,
		Avatar:   u.Avatar,
		IsAdmin:  u.IsAdmin,
	})
	return u.ID, nil
}
