package service

import (
	"github.com/isyayanoaa/go-blog/internal/dto"
	"github.com/isyayanoaa/go-blog/internal/model"
	"github.com/isyayanoaa/go-blog/internal/repo"
)

func ListPosts(category, tag string) ([]model.Post, error) {
	return repo.ListPosts(dto.PostFilter{Category: category, Tag: tag})
}
func GetPost(id string) (*model.Post, error) {
	post, err := repo.GetPostByID(id)
	if err != nil {
		return nil, err
	}
	go repo.IncrViewCount(id)
	return post, nil
}
func CreatePost(req dto.CreatePostReq) (uint, error)      { return repo.CreatePost(req) }
func UpdatePost(id string, req dto.UpdatePostReq) error   { return repo.UpdatePost(id, req) }
func DeletePost(id string) error                          { return repo.DeletePost(id) }
func ListCategories() ([]dto.CategoryResp, error)         { return repo.ListCategories() }
func ListTags() ([]dto.TagResp, error)                    { return repo.ListTags() }
func ListArchives() ([]dto.ArchiveResp, error)            { return repo.ListArchives() }
