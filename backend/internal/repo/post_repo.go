package repo

import (
	"strings"

	"github.com/isyayanoaa/go-blog/internal/dto"
	"github.com/isyayanoaa/go-blog/internal/model"
	"github.com/isyayanoaa/go-blog/pkg/database"
)

func ListPosts(filter dto.PostFilter) ([]model.Post, error) {
	var posts []model.Post
	q := database.DB.Where("status = ?", "published")
	if filter.Category != "" {
		q = q.Where("category = ?", filter.Category)
	}
	if filter.Tag != "" {
		q = q.Where("? = ANY(tags)", filter.Tag)
	}
	err := q.Order("created_at DESC").Find(&posts).Error
	return posts, err
}

func GetPostByID(id string) (*model.Post, error) {
	var post model.Post
	err := database.DB.First(&post, id).Error
	return &post, err
}

func IncrViewCount(id string) {
	database.DB.Model(&model.Post{}).Where("id = ?", id).
		UpdateColumn("view_count", database.DB.Raw("view_count + 1"))
}

func CreatePost(req dto.CreatePostReq) (uint, error) {
	if req.Status == "" {
		req.Status = "draft"
	}
	post := model.Post{
		Title:    req.Title,
		Slug:     strings.ToLower(strings.ReplaceAll(req.Title, " ", "-")),
		Content:  req.Content,
		Summary:  req.Summary,
		Cover:    req.Cover,
		Category: req.Category,
		Tags:     req.Tags,
		Status:   req.Status,
	}
	err := database.DB.Create(&post).Error
	return post.ID, err
}

func UpdatePost(id string, req dto.UpdatePostReq) error {
	return database.DB.Model(&model.Post{}).Where("id = ?", id).Updates(map[string]interface{}{
		"title": req.Title, "content": req.Content, "summary": req.Summary,
		"cover": req.Cover, "category": req.Category, "tags": req.Tags, "status": req.Status,
	}).Error
}

func DeletePost(id string) error {
	return database.DB.Delete(&model.Post{}, id).Error
}

func ListCategories() ([]dto.CategoryResp, error) {
	var cats []dto.CategoryResp
	err := database.DB.Model(&model.Post{}).
		Select("category as name, count(*) as count").
		Where("status = ? AND category != ?", "published", "").
		Group("category").Order("count DESC").Scan(&cats).Error
	return cats, err
}

func ListTags() ([]dto.TagResp, error) {
	var tags []dto.TagResp
	err := database.DB.Model(&model.Post{}).
		Select("unnest(tags) as name, count(*) as count").
		Where("status = ?", "published").
		Group("name").Order("count DESC").Scan(&tags).Error
	return tags, err
}

func ListArchives() ([]dto.ArchiveResp, error) {
	var arcs []dto.ArchiveResp
	err := database.DB.Model(&model.Post{}).
		Select("to_char(created_at, 'YYYY-MM') as month, count(*) as count").
		Where("status = ?", "published").
		Group("month").Order("month DESC").Scan(&arcs).Error
	return arcs, err
}
