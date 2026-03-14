package repo

import (
	"github.com/isyayanoaa/go-blog/internal/dto"
	"github.com/isyayanoaa/go-blog/internal/model"
	"github.com/isyayanoaa/go-blog/pkg/database"
)

func ListMoments() ([]model.Moment, error) {
	var moments []model.Moment
	err := database.DB.Order("created_at DESC").Find(&moments).Error
	return moments, err
}

func CreateMoment(req dto.CreateMomentReq) (uint, error) {
	m := model.Moment{Content: req.Content, Images: req.Images, Video: req.Video}
	err := database.DB.Create(&m).Error
	return m.ID, err
}

func DeleteMoment(id string) error {
	return database.DB.Delete(&model.Moment{}, id).Error
}
