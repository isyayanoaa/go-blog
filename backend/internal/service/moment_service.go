package service

import (
	"github.com/isyayanoaa/go-blog/internal/dto"
	"github.com/isyayanoaa/go-blog/internal/model"
	"github.com/isyayanoaa/go-blog/internal/repo"
)

func ListMoments() ([]model.Moment, error)                { return repo.ListMoments() }
func CreateMoment(req dto.CreateMomentReq) (uint, error)  { return repo.CreateMoment(req) }
func DeleteMoment(id string) error                        { return repo.DeleteMoment(id) }
